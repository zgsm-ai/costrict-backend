/**
 * @file number.ts
 * @description 数字格式化与精度安全工具
 */

const CREDIT_DECIMALS = 2;
const FLOAT_FIX_FACTOR = Math.pow(10, CREDIT_DECIMALS);

/**
 * 将数值四舍五入到 2 位小数，消除 0.1 + 0.2 = 0.30000000000000004 这类浮点误差。
 */
export function roundCredit(value: number): number {
    if (!Number.isFinite(value)) return 0;
    return Math.round(value * FLOAT_FIX_FACTOR) / FLOAT_FIX_FACTOR;
}

/**
 * 安全减法，结果经过 roundCredit 处理，避免浮点精度误差外泄。
 */
export function safeSubtract(a: number, b: number): number {
    if (!Number.isFinite(a) || !Number.isFinite(b)) return 0;
    return roundCredit(a - b);
}

/**
 * 格式化 Credit 数量用于展示：
 * - 整数直接显示（如 3500）
 * - 非整数保留 2 位小数（如 0.5、1.25）
 * - 内部先做 roundCredit 防止浮点噪声
 */
export function formatCredit(value: number): string {
    const rounded = roundCredit(value);
    return Number.isInteger(rounded) ? String(rounded) : rounded.toFixed(CREDIT_DECIMALS);
}

/**
 * 格式化 Credit 数量用于汇总展示，始终保留 2 位小数：
 * - 适用于「已使用 / 总额 / 剩余」这类需要稳定对齐的财务化展示
 * - 内部先做 roundCredit 防止浮点噪声
 */
export function formatCreditFixed(value: number): string {
    const rounded = roundCredit(value);
    return rounded.toFixed(CREDIT_DECIMALS);
}
