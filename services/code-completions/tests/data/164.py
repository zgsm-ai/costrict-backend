import pandas as pd
import matplotlib.pyplot as plt

def analyze_sales_data(file_path):
    df = pd.read_csv(file_path)
    df['date'] = pd.to_datetime(df['date'])
    
    <｜fim▁hole｜>monthly_sales = df.groupby(df['date'].dt.month)['amount'].sum()
    
    return monthly_sales

def plot_sales_trend(monthly_sales):
    plt.figure(figsize=(10, 6))
    monthly_sales.plot(kind='bar', color='skyblue')
    plt.title('Monthly Sales Trend')
    plt.xlabel('Month')
    plt.ylabel('Sales Amount')
    plt.savefig('sales_trend.png')
    plt.show()

def main():
    file_path = 'sales.csv'
    try:
        monthly_sales = analyze_sales_data(file_path)
        print(monthly_sales)
        plot_sales_trend(monthly_sales)
    except FileNotFoundError:
        print(f"Error: File '{file_path}' not found.")

if __name__ == "__main__":
    main()