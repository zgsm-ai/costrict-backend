import requests
import json

def fetch_user_data(user_id):
    url = f"https://api.example.com/users/{user_id}"
    
    try:
        response = requests.get(url)
        response.raise_for_status()
        
        <｜fim▁hole｜>user_data = response.json()
        
        return user_data
    except requests.exceptions.RequestException as e:
        print(f"Error fetching user data: {e}")
        return None

def main():
    user_id = 123
    user_data = fetch_user_data(user_id)
    
    if user_data:
        print(json.dumps(user_data, indent=2))
    else:
        print("Failed to fetch user data")

if __name__ == "__main__":
    main()