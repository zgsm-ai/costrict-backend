import numpy as np
from sklearn.model_selection import train_test_split
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_squared_error, r2_score

def generate_sample_data(n_samples=100):
    np.random.seed(42)
    X = np.random.rand(n_samples, 1) * 10
    y = 2 * X.squeeze() + 1 + np.random.randn(n_samples) * 2
    return X, y

def train_linear_regression(X, y):
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)
    
    <｜fim▁hole｜>model = LinearRegression()
    model.fit(X_train, y_train)
    
    y_test_pred = model.predict(X_test)
    test_mse = mean_squared_error(y_test, y_test_pred)
    test_r2 = r2_score(y_test, y_test_pred)
    
    return model, test_mse, test_r2

def main():
    X, y = generate_sample_data()
    model, test_mse, test_r2 = train_linear_regression(X, y)
    
    print(f"Test MSE: {test_mse:.4f}")
    print(f"Test R²: {test_r2:.4f}")

if __name__ == "__main__":
    main()