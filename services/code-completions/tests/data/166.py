import asyncio
import aiohttp
import time

async def fetch_url(session, url):
    try:
        start_time = time.time()
        async with session.get(url) as response:
            content = await response.text()
            end_time = time.time()
            
            return {
                'url': url,
                'status': response.status,
                'content_length': len(content),
                'response_time': end_time - start_time
            }
    except Exception as e:
        return {'url': url, 'error': str(e), 'response_time': 0}

async def fetch_all_urls(urls):
    <｜fim▁hole｜>async with aiohttp.ClientSession() as session:
        tasks = [fetch_url(session, url) for url in urls]
        return await asyncio.gather(*tasks)

async def main():
    urls = [
        "https://httpbin.org/get",
        "https://httpbin.org/delay/1",
        "https://httpbin.org/status/200"
    ]
    
    print(f"Fetching {len(urls)} URLs concurrently...")
    start_time = time.time()
    
    results = await fetch_all_urls(urls)
    
    end_time = time.time()
    total_time = end_time - start_time
    
    print("\nResults:")
    for result in results:
        if 'error' in result:
            print(f"{result['url']}: ERROR - {result['error']}")
        else:
            print(f"{result['url']}: Status {result['status']}, Size {result['content_length']} bytes, Time {result['response_time']:.3f}s")
    
    print(f"\nTotal time: {total_time:.3f} seconds")

if __name__ == "__main__":
    asyncio.run(main())