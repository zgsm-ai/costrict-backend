import datetime

class SocialMediaPlatform:
    def __init__(self):
        self.users = {}
        self.posts = {}
        self.comments = {}
        self.likes = {}
        self.friendships = {}
        self.notifications = {}
        self.trending_topics = []
    
    def create_user(self, user_id, username, email, password):
        if user_id in self.users:
            return False
        
        if any(user['email'] == email for user in self.users.values()):
            return False
        
        if any(user['username'] == username for user in self.users.values()):
            return False
        
        self.users[user_id] = {
            'username': username,
            'email': email,
            'password': password,
            'friends': [],
            'posts': [],
            'notifications': [],
            'created_at': self.get_current_timestamp()
        }
        return True
    
    def authenticate_user(self, email, password):
        for user_id, user_info in self.users.items():
            if user_info['email'] == email and user_info['password'] == password:
                return user_id
        return None
    
    def create_post(self, user_id, content, tags=None):
        if user_id not in self.users:
            return None
        
        post_id = self.generate_id('post')
        self.posts[post_id] = {
            'user_id': user_id,
            'content': content,
            'tags': tags or [],
            'likes': 0,
            'comments': [],
            'created_at': self.get_current_timestamp()
        }
        
        self.users[user_id]['posts'].append(post_id)
        return post_id
    
    def add_comment(self, user_id, post_id, content):
        if user_id not in self.users or post_id not in self.posts:
            return None
        
        comment_id = self.generate_id('comment')
        self.comments[comment_id] = {
            'user_id': user_id,
            'post_id': post_id,
            'content': content,
            'likes': 0,
            'created_at': self.get_current_timestamp()
        }
        
        self.posts[post_id]['comments'].append(comment_id)
        return comment_id
    
    def like_post(self, user_id, post_id):
        if user_id not in self.users or post_id not in self.posts:
            return False
        
        like_key = f"{user_id}_{post_id}"
        if like_key in self.likes:
            return False
        
        self.likes[like_key] = {
            'user_id': user_id,
            'post_id': post_id,
            'created_at': self.get_current_timestamp()
        }
        
        self.posts[post_id]['likes'] += 1
        return True
    
    def unlike_post(self, user_id, post_id):
        if user_id not in self.users or post_id not in self.posts:
            return False
        
        like_key = f"{user_id}_{post_id}"
        if like_key not in self.likes:
            return False
        
        del self.likes[like_key]
        self.posts[post_id]['likes'] -= 1
        return True
    
    def send_friend_request(self, sender_id, receiver_id):
        if sender_id not in self.users or receiver_id not in self.users:
            return False
        
        if sender_id == receiver_id:
            return False
        
        friendship_key = f"{sender_id}_{receiver_id}"
        if friendship_key in self.friendships:
            return False
        
        self.friendships[friendship_key] = {
            'user1_id': sender_id,
            'user2_id': receiver_id,
            'status': 'pending',
            'created_at': self.get_current_timestamp()
        }
        return True
    
    def accept_friend_request(self, sender_id, receiver_id):
        friendship_key = f"{sender_id}_{receiver_id}"
        if friendship_key not in self.friendships:
            return False
        
        if self.friendships[friendship_key]['status'] != 'pending':
            return False
        
        self.friendships[friendship_key]['status'] = 'accepted'
        self.friendships[friendship_key]['updated_at'] = self.get_current_timestamp()
        
        self.users[sender_id]['friends'].append(receiver_id)
        self.users[receiver_id]['friends'].append(sender_id)
        
        return True
    
    def get_user_feed(self, user_id, limit=10):
        if user_id not in self.users:
            return []
        
        user_and_friends = [user_id] + self.users[user_id]['friends']
        
        feed_posts = []
        for post_id, post_info in self.posts.items():
            if post_info['user_id'] in user_and_friends:
                post_copy = post_info.copy()
                post_copy['post_id'] = post_id
                post_copy['author'] = self.users[post_info['user_id']]['username']
                feed_posts.append(post_copy)
        
        feed_posts.sort(key=lambda x: x['created_at'], reverse=True)
        
        if len(feed_posts) > limit:
            feed_posts = feed_posts[:limit]
        
        return feed_posts
    
    def search_posts(self, keyword):
        keyword = keyword.lower()
        result = []
        
        for post_id, post_info in self.posts.items():
            if (keyword in post_info['content'].lower() or
                any(keyword in tag.lower() for tag in post_info['tags'])):
                post_copy = post_info.copy()
                post_copy['post_id'] = post_id
                post_copy['author'] = self.users[post_info['user_id']]['username']
                result.append(post_copy)
        
        return result
    
    def generate_id(self, prefix):
        import uuid
        return f"{prefix}_{uuid.uuid4().hex[:8]}"
    
    def get_current_timestamp(self):
        return datetime.datetime.now().isoformat()

# 示例使用
if __name__ == "__main__":
    platform = SocialMediaPlatform()
    
    # 创建用户
    user1_id = "user1"
    platform.create_user(user1_id, "alice", "alice@example.com", "password123")
    
    user2_id = "user2"
    platform.create_user(user2_id, "bob", "bob@example.com", "password456")
    
    # 创建帖子
    post1_id = platform.create_post(user1_id, "Hello, this is my first post!", ["hello", "first"])
    post2_id = platform.create_post(user2_id, "Nice to meet you all!", ["greeting"])
    
    # 点赞帖子
    platform.like_post(user2_id, post1_id)
    
    # 添加评论
    platform.add_comment(user2_id, post1_id, "Great post!")
    
    # 发送好友请求
    platform.send_friend_request(user1_id, user2_id)
    
    # 接受好友请求
    platform.accept_friend_request(user1_id, user2_id)
    
    # 获取用户动态
    feed = platform.get_user_feed(user1_id)
    print("User Feed:")
    for post in feed:
        print(f"- {post['author']}: {post['content']}")
    
    # 搜索帖子
    search_results = platform.search_posts("hello")
    print("\nSearch Results:")
    for post in search_results:
        print(f"- {post['author']}: {post['content']}")