/**
 * @file order mock
 */

import type { MockMethod, Recordable } from 'vite-plugin-mock';
import { successResponse } from './_utils';

const mockOrders = [
    {
        id: 1,
        order_id: 'ORD20250501001',
        user_id: 'user-mock-001',
        amount: 9.9,
        status: 'completed',
        quota_type: '探索版',
        order_source: 'web',
        credit_count: 1000,
        credit_expire_date: '2026-06-01',
        created_at: '2026-05-01T10:00:00Z',
        updated_at: '2026-05-01T10:00:00Z',
        invoice_status: 1,
        description: '探索版套餐',
        jd_order_id: '',
    },
    {
        id: 2,
        order_id: 'ORD20250420002',
        user_id: 'user-mock-001',
        amount: 98,
        status: 'completed',
        quota_type: '专业版',
        order_source: 'web',
        credit_count: 4200,
        credit_expire_date: '2026-05-20',
        created_at: '2026-04-20T14:30:00Z',
        updated_at: '2026-04-20T14:30:00Z',
        invoice_status: 0,
        description: '专业版套餐',
        jd_order_id: '',
    },
    {
        id: 3,
        order_id: 'ORD20250410003',
        user_id: 'user-mock-001',
        amount: 248,
        status: 'completed',
        quota_type: '旗舰版',
        order_source: 'web',
        credit_count: 10800,
        credit_expire_date: '2026-05-10',
        created_at: '2026-04-10T09:00:00Z',
        updated_at: '2026-04-10T09:00:00Z',
        invoice_status: 2,
        description: '旗舰版套餐',
        jd_order_id: '',
    },
];

const mockQuotaTypes = [
    {
        id: 4,
        display_name: '探索版',
        credit_count: 1000,
        amount: 9.9,
        valid_days: 30,
        quota_marketing_rules_id: 1,
        original_amount: 50,
        equivalent_credits: 1000,
        bonus_credits: 0,
        estimated_requests: 1000,
        created_at: '2025-12-06T11:03:51.952158Z',
        updated_at: '2025-12-06T11:03:51.952158Z',
        marketing_rules: {
            id: 1,
            rule_type: 1,
            pay_amount: 9.9,
            pay_discount: 1,
            status: 1,
            apply_user_type: 1,
            created_at: '2025-12-06T10:18:37.72414Z',
            updated_at: '2025-12-06T10:18:37.72414Z',
        },
    },
    {
        id: 5,
        display_name: '专业版',
        credit_count: 4200,
        amount: 98,
        valid_days: 30,
        quota_marketing_rules_id: 0,
        original_amount: 200,
        equivalent_credits: 4000,
        bonus_credits: 200,
        estimated_requests: 4200,
        created_at: '2025-12-06T11:03:51.952158Z',
        updated_at: '2025-12-06T11:03:51.952158Z',
    },
    {
        id: 6,
        display_name: '旗舰版',
        credit_count: 10800,
        amount: 248,
        valid_days: 30,
        quota_marketing_rules_id: 0,
        original_amount: 500,
        equivalent_credits: 10000,
        bonus_credits: 800,
        estimated_requests: 10800,
        created_at: '2025-12-06T11:03:51.952158Z',
        updated_at: '2025-12-06T11:03:51.952158Z',
    },
];

export default [
    {
        url: '/quota-order-manager/api/v1/orders',
        method: 'get',
        response: ({ query }: { query: Recordable }) => {
            const page = Number(query.page) || 1;
            const pageSize = Number(query.page_size) || 10;
            const total = mockOrders.length;
            const start = (page - 1) * pageSize;
            const end = start + pageSize;
            return successResponse({
                orders: mockOrders.slice(start, end),
                total,
                limit: pageSize,
                offset: (page - 1) * pageSize,
            });
        },
    },
    {
        url: '/quota-order-manager/api/v1/orders',
        method: 'post',
        response: () =>
            successResponse({
                id: 4,
                order_id: 'ORD20250507004',
                user_id: 'user-mock-001',
                amount: 9.9,
                credit_count: 1000,
                credit_expire_date: '2026-06-07',
                status: 'pending',
                quota_type: '探索版',
                description: '探索版套餐',
                order_source: 'web',
                jd_order_id: '',
                invoice_status: 0,
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString(),
            }),
    },
    {
        url: '/quota-order-manager/api/v1/payment/initiate',
        method: 'post',
        response: () =>
            successResponse({
                version: '1.0',
                encoding: 'UTF-8',
                sign: 'mock_sign_xxx',
                signMethod: 'RSA',
                returnCode: 'SUCCESS',
                respCode: '00',
                qrCode: 'weixin://wxpay/bizpayurl?pr=mock123',
                cmbOrderId: 'CMB20250507001',
                txnTime: '20260507100000',
                biz_content: {
                    merId: 'M001',
                    orderId: 'ORD20250507004',
                    cmbOrderId: 'CMB20250507001',
                    qrCode: 'weixin://wxpay/bizpayurl?pr=mock123',
                    txnTime: '20260507100000',
                },
            }),
    },
    {
        url: '/quota-order-manager/api/v1/invoices',
        method: 'post',
        response: () =>
            successResponse({
                invoice_id: 1,
                title_type: 1,
                invoice_type: 1,
                invoice_title: 'Mock Company Ltd.',
                taxpayer_id: '91110000123456789X',
                company_address: 'Beijing, China',
                company_phone: '',
                bank_name: 'Industrial and Commercial Bank of China',
                bank_account: '6222021234567890123',
                receive_email: 'invoice@example.com',
                order_id: 'ORD20250501001',
                amount: 9.9,
                invoice_content: 'Software Service',
                invoice_status: 0,
                apply_time: new Date().toISOString(),
                issue_time: null,
            }),
    },
    {
        url: '/quota-order-manager/api/v1/orders/:id',
        method: 'get',
        response: ({ query }: { query: Recordable }) => {
            const orderId = query.id as string;
            const order = mockOrders.find((o) => o.order_id === orderId) || mockOrders[0];
            return successResponse(order);
        },
    },
    {
        url: '/quota-order-manager/api/v1/quotas/types',
        method: 'get',
        response: () => successResponse(mockQuotaTypes),
    },
    {
        url: '/quota-order-manager/api/v1/quotas/types/:id',
        method: 'get',
        response: ({ query }: { query: Recordable }) => {
            const id = Number(query.id);
            const type = mockQuotaTypes.find((t) => t.id === id) || mockQuotaTypes[0];
            const quantity = Number(query.quantity) || 1;
            return successResponse({
                ...type,
                quantity,
                amount: type.amount * quantity,
            });
        },
    },
] as MockMethod[];
