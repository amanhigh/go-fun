"""
Stock Price Analyzer for Schedule Foreign Assets (FA) in Indian ITR.

Description:
This script analyzes stock prices for specified foreign assets (Schedule FA) 
and calculates their values in Indian Rupees (INR) using SBI's TT Buy rates. 
It provides peak and year-end prices for a given tax year, aiding in Indian 
Income Tax Return (ITR) preparation.

https://zerodha.com/z-connect/varsity/how-is-the-income-from-foreign-equity-shares-taxed-and-disclosed-in-the-itr-form

Features:
- Downloads stock data from Alpha Vantage API
- Uses SBI Reference Rates for USD-INR conversion
- Calculates peak and year-end prices for specified stocks
- Converts USD prices to INR using relevant exchange rates
- Displays results in a formatted table
- Caches downloaded data to avoid unnecessary API calls
- Uses TT Buy rate for currency conversion

Requirements:
- Python 3.6+
- Required packages: requests, json, csv, datetime

Setup:
1. Install required packages:
   pip install requests

2. Set up an environment variable for your Alpha Vantage API key:
   export VANTAGE_API_KEY='your_api_key_here'

Configuration:
Edit the following constants in the script as needed:
- TICKERS: List of stock tickers to analyze
- DOWNLOADS_DIR: Directory to store downloaded data
- TAX_YEAR: The year for which to perform the analysis

Usage:
Run the script from the command line:
python schedule_fa.py

Output:
The script displays a table with the following information for each ticker:
- Ticker symbol
- Peak date and price (USD)
- Year-end date and price (USD)
- TT Buy Rate for peak and year-end dates
- Peak and year-end prices converted to INR

Example Output:
python schedule_fa.py
Tickers: ['AMZN', 'SIVR']
Tax Year: 2023
API key loaded successfully
Directory created/verified: ~/Downloads/Tickers
Data for AMZN already exists (last modified: 2024-09-07). Skipping download.
Data for SIVR already exists (last modified: 2024-09-07). Skipping download.
SBI USD rates data already exists (last modified: 2024-09-07). Skipping download.

Ticker Data:
Ticker | Peak Date  | Peak Price (USD) | Year-End Date | Year-End Price (USD) | TTBR (Peak) | TTBR (Year-End) | Peak Price (INR) | Year-End Price (INR)
-------------------------------------------------------------------------------------------------------------------------------------------------------
AMZN   | 2023-12-18 | $154.07          | 2023-12-29    | $151.94              | ₹82.57      | ₹82.76          | ₹12721.56        | ₹12574.55           
SIVR   | 2023-05-04 | $25.00           | 2023-12-29    | $22.77               | ₹81.27      | ₹82.76          | ₹2031.75         | ₹1884.45

Troubleshooting:
- "API key not found" error: Ensure the VANTAGE_API_KEY environment variable is set.
- "File not found" errors: Check if the DOWNLOADS_DIR path is correct and accessible.

Disclaimer:
This is a helper script intended for informational purposes only. The author 
assumes no legal or financial responsibility for the usage of this script. 
Users should consult their Chartered Accountant (CA) for validation and 
professional advice regarding their tax filings and financial matters. The 
information provided by this script should not be considered as financial 
or legal advice.
"""

import os
import requests
import json
import csv
from datetime import datetime
# TODO: #A Migrate Completely
# FIXME: #C Create Readme on how to use FA.
# Constants
TICKERS = ['AMZN', 'SIVR']
DOWNLOADS_DIR = '~/Downloads/Tickers'
BASE_URL = 'https://www.alphavantage.co/query'
API_KEY_ENV_VAR = 'VANTAGE_API_KEY'
TAX_YEAR = 2023
SBI_USD_RATES_URL = 'https://raw.githubusercontent.com/sahilgupta/sbi-fx-ratekeeper/main/csv_files/SBI_REFERENCE_RATES_USD.csv'

# Moved
def get_api_key():
    api_key = os.getenv(API_KEY_ENV_VAR)
    if not api_key:
        raise ValueError(f"API key not found. Please set the '{API_KEY_ENV_VAR}' environment variable.")
    return api_key

def create_download_directory(dir_path):
    expanded_path = os.path.expanduser(dir_path)
    os.makedirs(expanded_path, exist_ok=True)
    print(f"Directory created/verified: {expanded_path}")
    return expanded_path

def check_file_exists(file_path):
    if os.path.exists(file_path):
        mod_time = os.path.getmtime(file_path)
        mod_date = datetime.fromtimestamp(mod_time).strftime('%Y-%m-%d')
        return True, mod_date
    return False, None

def download_ticker_data(ticker, api_key):
    params = {'function': 'TIME_SERIES_DAILY', 'symbol': ticker, 'outputsize': 'full', 'apikey': api_key}
    response = requests.get(BASE_URL, params=params)
    return response.json()

def save_ticker_data(data, file_path):
    with open(file_path, 'w') as f:
        json.dump(data, f)

# // FIXME: #B Migrate to Ticker Manager as DownloadTicker
def download_tickers(tickers, downloads_dir, api_key):
    for ticker in tickers:
        file_path = os.path.join(downloads_dir, f'{ticker}.json')
        file_exists, mod_date = check_file_exists(file_path)
        
        if file_exists:
            print(f'Data for {ticker} already exists (last modified: {mod_date}). Skipping download.')
        else:
            data = download_ticker_data(ticker, api_key)
            save_ticker_data(data, file_path)
            print(f'Data for {ticker} downloaded and saved to {file_path}')

# Moved
def download_sbi_usd_rates(downloads_dir):
    file_path = os.path.join(downloads_dir, 'SBI_REFERENCE_RATES_USD.csv')
    file_exists, mod_date = check_file_exists(file_path)
    
    if file_exists:
        print(f'SBI USD rates data already exists (last modified: {mod_date}). Skipping download.')
    else:
        response = requests.get(SBI_USD_RATES_URL)
        with open(file_path, 'wb') as f:
            f.write(response.content)
        print(f'SBI USD rates data downloaded and saved to {file_path}')

# // FIXME: #B Migrate to Ticker Manager as AnalyseTicker
def find_ticker_data(ticker, downloads_dir, year):
    file_path = os.path.join(downloads_dir, f'{ticker}.json')
    with open(file_path, 'r') as f:
        data = json.load(f)
    
    time_series = data['Time Series (Daily)']
    highest_close = 0
    highest_date = None
    year_end_close = None
    year_end_date = f"{year}-12-31"
    
    for date, values in time_series.items():
        if date.startswith(str(year)):
            close_price = float(values['4. close'])
            if close_price > highest_close:
                highest_close = close_price
                highest_date = date
            if date == year_end_date:
                year_end_close = close_price
    
    if year_end_close is None:
        last_trading_day = max(date for date in time_series.keys() if date.startswith(str(year)))
        year_end_close = float(time_series[last_trading_day]['4. close'])
        year_end_date = last_trading_day

    # Get TTBR for peak and year-end dates
    peak_ttbr = find_sbi_usd_rate(downloads_dir, highest_date)
    year_end_ttbr = find_sbi_usd_rate(downloads_dir, year_end_date)

    return highest_date, highest_close, year_end_date, year_end_close, peak_ttbr, year_end_ttbr

#  FIXME: #A Migrate to SBI Manager
def find_sbi_usd_rate(downloads_dir, date):
    file_path = os.path.join(downloads_dir, 'SBI_REFERENCE_RATES_USD.csv')
    with open(file_path, 'r') as f:
        csv_reader = csv.DictReader(f)
        for row in csv_reader:
            if row['DATE'].split()[0] == date:
                return float(row['TT BUY'])
    return None

# // TODO: Migrate to FAManager
def process_tickers(tickers, downloads_dir, year):
    table_data = []
    for ticker in tickers:
        highest_date, highest_close, year_end_date, year_end_close, peak_ttbr, year_end_ttbr = find_ticker_data(ticker, downloads_dir, year)
        if highest_date:
            table_data.append([
                ticker,
                highest_date,
                f"${highest_close:.2f}",
                year_end_date,
                f"${year_end_close:.2f}",
                f"₹{peak_ttbr:.2f}" if peak_ttbr else "N/A",
                f"₹{year_end_ttbr:.2f}" if year_end_ttbr else "N/A",
                f"₹{highest_close * peak_ttbr:.2f}" if peak_ttbr else "N/A",
                f"₹{year_end_close * year_end_ttbr:.2f}" if year_end_ttbr else "N/A"
            ])
        else:
            table_data.append([ticker, "No data", "No data", "No data", "No data", "No data", "No data", "No data", "No data"])
    return table_data

def print_table(data):
    headers = ["Ticker", "Peak Date", "Peak Price (USD)", "Year-End Date", "Year-End Price (USD)", "TTBR (Peak)", "TTBR (Year-End)", "Peak Price (INR)", "Year-End Price (INR)"]
    col_widths = [max(len(str(row[i])) for row in data + [headers]) for i in range(len(headers))]
    
    header_row = " | ".join(f"{headers[i]:<{col_widths[i]}}" for i in range(len(headers)))
    print(header_row)
    print("-" * len(header_row))
    
    for row in data:
        print(" | ".join(f"{str(row[i]):<{col_widths[i]}}" for i in range(len(row))))

def main():
    print("Tickers:", TICKERS)
    print(f"Tax Year: {TAX_YEAR}")

    api_key = get_api_key()
    print("API key loaded successfully")

    downloads_dir = create_download_directory(DOWNLOADS_DIR)

    download_tickers(TICKERS, downloads_dir, api_key)
    download_sbi_usd_rates(downloads_dir)

    table_data = process_tickers(TICKERS, downloads_dir, TAX_YEAR)

    print("\nTicker Data:")
    print_table(table_data)

    print("\nAll data processed successfully")

if __name__ == "__main__":
    main()
