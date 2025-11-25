async function fetchData(url: string): Promise<any> {
    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching data:', error);
        return null;
    }
}

async function main() {
    const apiUrl = 'https://api.example.com/users/1';
    <｜fim▁hole｜>const userData = await fetchData(apiUrl);
    if (userData) {
        console.log('User data:', userData);
    } else {
        console.log('Failed to fetch user data');
    }
}

main();