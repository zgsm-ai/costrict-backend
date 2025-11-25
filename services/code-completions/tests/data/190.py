class InventoryManagementSystem:
    """
    库存管理系统类，用于管理产品的库存、采购和销售
    """
    
    def __init__(self):
        self.products = {}  # 产品ID到产品信息的映射
        self.suppliers = {}  # 供应商ID到供应商信息的映射
        self.transactions = []  # 交易记录
        self.reports = []  # 报告记录
    
    def add_product(self, product_id, name, category, price, quantity, supplier_id):
        """
        添加新产品到库存系统
        @param product_id: 产品唯一标识符
        @param name: 产品名称
        @param category: 产品类别
        @param price: 产品价格
        @param quantity: 初始库存数量
        @param supplier_id: 供应商ID
        @return: 添加成功返回True，否则返回False
        """
        if product_id in self.products:
            return False
        
        if supplier_id not in self.suppliers:
            return False
        
        self.products[product_id] = {
            'name': name,
            'category': category,
            'price': price,
            'quantity': quantity,
            'supplier_id': supplier_id
        }
        return True
    
    def add_supplier(self, supplier_id, name, contact_person, email, phone):
        """
        添加新供应商到系统
        @param supplier_id: 供应商唯一标识符
        @param name: 供应商名称
        @param contact_person: 联系人姓名
        @param email: 电子邮件地址
        @param phone: 联系电话
        @return: 添加成功返回True，否则返回False
        """
        if supplier_id in self.suppliers:
            return False
        
        self.suppliers[supplier_id] = {
            'name': name,
            'contact_person': contact_person,
            'email': email,
            'phone': phone
        }
        return True
    
    def update_stock(self, product_id, quantity_change, transaction_type, notes=""):
        """
        更新产品库存数量并记录交易
        @param product_id: 产品ID
        @param quantity_change: 数量变化（正值为增加，负值为减少）
        @param transaction_type: 交易类型（如"purchase", "sale", "return"）
        @param notes: 交易备注信息
        @return: 更新成功返回True，否则返回False
        """
        if product_id not in self.products:
            return False
        
        new_quantity = self.products[product_id]['quantity'] + quantity_change
        
        if new_quantity < 0:
            return False
        
        self.products[product_id]['quantity'] = new_quantity
        
        # 记录交易
        transaction = {
            'product_id': product_id,
            'quantity_change': quantity_change,
            'transaction_type': transaction_type,
            'timestamp': self.get_current_timestamp(),
            'notes': notes
        }
        self.transactions.append(transaction)
        return True
    
    def get_product_info(self, product_id):
        """
        获取产品信息
        @param product_id: 产品ID
        @return: 产品信息字典，如果不存在则返回None
        """
        if product_id not in self.products:
            return None
        
        product_info = self.products[product_id].copy()
        supplier_id = product_info['supplier_id']
        if supplier_id in self.suppliers:
            product_info['supplier_info'] = self.suppliers[supplier_id].copy()
        
        return product_info
    
    def get_supplier_info(self, supplier_id):
        """
        获取供应商信息
        @param supplier_id: 供应商ID
        @return: 供应商信息字典，如果不存在则返回None
        """
        if supplier_id not in self.suppliers:
            return None
        
        return self.suppliers[supplier_id].copy()
    
    def get_products_by_supplier(self, supplier_id):
        """
        获取指定供应商的所有产品
        @param supplier_id: 供应商ID
        @return: 产品列表
        """
        if supplier_id not in self.suppliers:
            return []
        
        result = []
        for product_id, product_info in self.products.items():
            if product_info['supplier_id'] == supplier_id:
                product_copy = product_info.copy()
                product_copy['product_id'] = product_id
                result.append(product_copy)
        
        return result
    
    def get_products_by_category(self, category):
        """
        获取指定类别的所有产品
        @param category: 产品类别
        @return: 产品列表
        """
        result = []
        for product_id, product_info in self.products.items():
            if product_info['category'] == category:
                product_copy = product_info.copy()
                product_copy['product_id'] = product_id
                
                # 添加供应商信息
                supplier_id = product_info['supplier_id']
                if supplier_id in self.suppliers:
                    product_copy['supplier_info'] = self.suppliers[supplier_id].copy()
                
                result.append(product_copy)
        
        return result
    
    def get_low_stock_products(self, threshold=10):
        """
        获取库存低于阈值的产品
        @param threshold: 库存阈值
        @return: 低库存产品列表
        """
        result = []
        for product_id, product_info in self.products.items():
            if product_info['quantity'] < threshold:
                product_copy = product_info.copy()
                product_copy['product_id'] = product_id
                result.append(product_copy)
        
        return result
    
    def generate_inventory_report(self):
        """
        生成库存报告
        @return: 报告字典
        """
        total_products = len(self.products)
        total_suppliers = len(self.suppliers)
        total_transactions = len(self.transactions)
        
        total_value = 0
        for product_info in self.products.values():
            total_value += product_info['price'] * product_info['quantity']
        
        category_stats = {}
        for product_info in self.products.values():
            category = product_info['category']
            if category not in category_stats:
                category_stats[category] = {
                    'count': 0,
                    'total_value': 0
                }
            category_stats[category]['count'] += 1
            category_stats[category]['total_value'] += product_info['price'] * product_info['quantity']
        
        report = {
            'timestamp': self.get_current_timestamp(),
            'total_products': total_products,
            'total_suppliers': total_suppliers,
            'total_transactions': total_transactions,
            'total_inventory_value': total_value,
            'category_statistics': category_stats,
            'low_stock_products': self.get_low_stock_products()
        }
        
        self.reports.append(report)
        return report
    
    def search_products(self, keyword):
        """
        根据关键词搜索产品
        @param keyword: 搜索关键词
        @return: 匹配的产品列表
        """
        keyword = keyword.lower()
        result = []
        
        for product_id, product_info in self.products.items():
            if (keyword in product_info['name'].lower() or 
                keyword in product_info['category'].lower()):
                product_copy = product_info.copy()
                product_copy['product_id'] = product_id
                result.append(product_copy)
        
        return result
    
    def get_transaction_history(self, product_id=None, transaction_type=None, limit=None):
        """
        获取交易历史记录
        @param product_id: 按产品ID筛选（可选）
        @param transaction_type: 按交易类型筛选（可选）
        @param limit: 返回记录数限制（可选）
        @return: 交易记录列表
        """
        result = []
        
        for transaction in self.transactions:
            if product_id is not None and transaction['product_id'] != product_id:
                continue
            
            if transaction_type is not None and transaction['transaction_type'] != transaction_type:
                continue
            
            result.append(transaction.copy())
        
        # 按时间戳降序排序
        result.sort(key=lambda x: x['timestamp'], reverse=True)
        
        if limit is not None and len(result) > limit:
            result = result[:limit]
        
        return result
    
    def get_current_timestamp(self):
        """
        获取当前时间戳
        @return: 当前时间戳字符串
        """
        import datetime
        return datetime.datetime.now().isoformat()