package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

type Tensor struct {
	Data  []float64
	Shape []int
	Grad  *Tensor
}

func NewTensor(data []float64, shape []int) *Tensor {
	return &Tensor{
		Data:  data,
		Shape: shape,
	}
}

func (t *Tensor) Size() int {
	size := 1
	for _, dim := range t.Shape {
		size *= dim
	}
	return size
}

func (t *Tensor) ZeroGrad() {
	if t.Grad != nil {
		for i := range t.Grad.Data {
			t.Grad.Data[i] = 0
		}
	}
}

type Optimizer interface {
	Step(params map[string]*Tensor)
	ZeroGrad()
	SetLR(float64)
}

type SGD struct {
	lr     float64
	params map[string]*Tensor
}

func NewSGD(params map[string]*Tensor, lr float64) *SGD {
	return &SGD{
		lr:     lr,
		params: params,
	}
}

func (sgd *SGD) Step(params map[string]*Tensor) {
	for _, param := range params {
		if param.Grad == nil {
			continue
		}
		for i := 0; i < param.Size(); i++ {
			param.Data[i] -= sgd.lr * param.Grad.Data[i]
		}
	}
}

func (sgd *SGD) ZeroGrad() {
	for _, param := range sgd.params {
		param.ZeroGrad()
	}
}

func (sgd *SGD) SetLR(lr float64) {
	sgd.lr = lr
}

type Layer interface {
	Forward(input *Tensor) *Tensor
	Backward(gradOutput *Tensor) *Tensor
	GetParameters() map[string]*Tensor
}

type Linear struct {
	weights *Tensor
	bias    *Tensor
	input   *Tensor
	output  *Tensor
}

func NewLinear(inFeatures, outFeatures int, bias bool) *Linear {
	weightData := make([]float64, inFeatures*outFeatures)
	for i := range weightData {
		weightData[i] = rand.NormFloat64() * math.Sqrt(2.0/float64(inFeatures))
	}
	weights := NewTensor(weightData, []int{outFeatures, inFeatures})
	weights.Grad = NewTensor(make([]float64, weights.Size()), weights.Shape)

	var biasTensor *Tensor
	if bias {
		biasData := make([]float64, outFeatures)
		biasTensor = NewTensor(biasData, []int{outFeatures})
		biasTensor.Grad = NewTensor(make([]float64, biasTensor.Size()), biasTensor.Shape)
	}

	return &Linear{
		weights: weights,
		bias:    biasTensor,
	}
}

func (l *Linear) Forward(input *Tensor) *Tensor {
	l.input = input
	output := l.matmul(l.weights, input)
	if l.bias != nil {
		output = l.addBias(output, l.bias)
	}
	l.output = output
	return output
}

func (l *Linear) Backward(gradOutput *Tensor) *Tensor {
	if l.weights.Grad != nil {
		gradWeights := l.matmul(gradOutput, l.transpose(l.input))
		for i := range gradWeights.Data {
			l.weights.Grad.Data[i] += gradWeights.Data[i]
		}
	}

	if l.bias != nil && l.bias.Grad != nil {
		gradBias := l.sumOverBatch(gradOutput)
		for i := range gradBias.Data {
			l.bias.Grad.Data[i] += gradBias.Data[i]
		}
	}

	return l.matmul(l.transpose(l.weights), gradOutput)
}

func (l *Linear) GetParameters() map[string]*Tensor {
	params := make(map[string]*Tensor)
	params["weights"] = l.weights
	if l.bias != nil {
		params["bias"] = l.bias
	}
	return params
}

func (l *Linear) matmul(a, b *Tensor) *Tensor {
	m := a.Shape[0]
	n := b.Shape[1]
	k := a.Shape[1]

	result := make([]float64, m*n)

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			sum := 0.0
			for p := 0; p < k; p++ {
				sum += a.Data[i*k+p] * b.Data[p*n+j]
			}
			result[i*n+j] = sum
		}
	}

	return NewTensor(result, []int{m, n})
}

func (l *Linear) transpose(a *Tensor) *Tensor {
	m := a.Shape[0]
	n := a.Shape[1]

	result := make([]float64, m*n)

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			result[j*m+i] = a.Data[i*n+j]
		}
	}

	return NewTensor(result, []int{n, m})
}

func (l *Linear) addBias(a, bias *Tensor) *Tensor {
	m := a.Shape[0]
	n := a.Shape[1]

	result := make([]float64, m*n)

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			result[i*n+j] = a.Data[i*n+j] + bias.Data[j]
		}
	}

	return NewTensor(result, a.Shape)
}

func (l *Linear) sumOverBatch(a *Tensor) *Tensor {
	m := a.Shape[0]
	n := a.Shape[1]

	result := make([]float64, n)

	for j := 0; j < n; j++ {
		sum := 0.0
		for i := 0; i < m; i++ {
			sum += a.Data[i*n+j]
		}
		result[j] = sum
	}

	return NewTensor(result, []int{n})
}

type ReLU struct {
	input  *Tensor
	output *Tensor
}

func NewReLU() *ReLU {
	return &ReLU{}
}

func (r *ReLU) Forward(input *Tensor) *Tensor {
	r.input = input
	result := make([]float64, input.Size())
	for i, val := range input.Data {
		if val > 0 {
			result[i] = val
		} else {
			result[i] = 0
		}
	}
	r.output = NewTensor(result, input.Shape)
	return r.output
}

func (r *ReLU) Backward(gradOutput *Tensor) *Tensor {
	result := make([]float64, r.input.Size())
	for i, val := range r.input.Data {
		if val > 0 {
			result[i] = gradOutput.Data[i]
		} else {
			result[i] = 0
		}
	}
	return NewTensor(result, r.input.Shape)
}

func (r *ReLU) GetParameters() map[string]*Tensor {
	return make(map[string]*Tensor)
}

type MSELoss struct {
	predict *Tensor
	target  *Tensor
	output  float64
}

func NewMSELoss() *MSELoss {
	return &MSELoss{}
}

func (l *MSELoss) Forward(predict, target *Tensor) float64 {
	l.predict = predict
	l.target = target

	sum := 0.0
	for i := 0; i < predict.Size(); i++ {
		diff := predict.Data[i] - target.Data[i]
		sum += diff * diff
	}

	l.output = sum / float64(predict.Size())
	return l.output
}

func (l *MSELoss) Backward() *Tensor {
	result := make([]float64, l.predict.Size())
	for i := 0; i < l.predict.Size(); i++ {
		result[i] = 2.0 * (l.predict.Data[i] - l.target.Data[i]) / float64(l.predict.Size())
	}
	return NewTensor(result, l.predict.Shape)
}

type Model struct {
	layers []Layer
}

func NewModel(layers []Layer) *Model {
	return &Model{
		layers: layers,
	}
}

func (m *Model) Forward(input *Tensor) *Tensor {
	output := input
	for _, layer := range m.layers {
		output = layer.Forward(output)
	}
	return output
}

func (m *Model) Backward(gradOutput *Tensor) *Tensor {
	grad := gradOutput
	for i := len(m.layers) - 1; i >= 0; i-- {
		grad = m.layers[i].Backward(grad)
	}
	return grad
}

func (m *Model) GetParameters() map[string]*Tensor {
	params := make(map[string]*Tensor)
	for i, layer := range m.layers {
		layerParams := layer.GetParameters()
		for name, param := range layerParams {
			params[fmt.Sprintf("layer_%d_%s", i, name)] = param
		}
	}
	return params
}

type Trainer struct {
	model     *Model
	optimizer Optimizer
	loss      *MSELoss
}

func NewTrainer(model *Model, optimizer Optimizer) *Trainer {
	return &Trainer{
		model:     model,
		optimizer: optimizer,
		loss:      NewMSELoss(),
	}
}

func (t *Trainer) TrainStep(input, target *Tensor) float64 {
	output := t.model.Forward(input)
	loss := t.loss.Forward(output, target)
	gradOutput := t.loss.Backward()
	t.model.Backward(gradOutput)
	t.optimizer.Step(t.model.GetParameters())
	t.optimizer.ZeroGrad()
	return loss
}

func main() {
	rand.Seed(time.Now().UnixNano())

	inputSize := 10
	hiddenSize := 5
	outputSize := 1
	<｜fim▁hole｜>

	model := NewModel([]Layer{
		NewLinear(inputSize, hiddenSize, true),
		NewReLU(),
		NewLinear(hiddenSize, outputSize, true),
	})

	params := model.GetParameters()
	optimizer := NewSGD(params, 0.01)

	trainer := NewTrainer(model, optimizer)

	numSamples := 100
	data := make([]*Tensor, numSamples)
	targets := make([]*Tensor, numSamples)

	for i := 0; i < numSamples; i++ {
		inputData := make([]float64, inputSize)
		for j := 0; j < inputSize; j++ {
			inputData[j] = rand.NormFloat64()
		}
		data[i] = NewTensor(inputData, []int{inputSize})

		targetData := make([]float64, outputSize)
		for j := 0; j < outputSize; j++ {
			targetData[j] = rand.NormFloat64()
		}
		targets[i] = NewTensor(targetData, []int{outputSize})
	}

	epochs := 10
	for epoch := 1; epoch <= epochs; epoch++ {
		totalLoss := 0.0
		for i := 0; i < numSamples; i++ {
			loss := trainer.TrainStep(data[i], targets[i])
			totalLoss += loss
		}
		avgLoss := totalLoss / float64(numSamples)
		log.Printf("Epoch %d/%d, Loss: %.4f", epoch, epochs, avgLoss)
	}

	testInput := NewTensor([]float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0}, []int{inputSize})
	prediction := model.Forward(testInput)
	fmt.Printf("预测结果: %v\n", prediction.Data)
	fmt.Println("完成")
}
