import Cookies from 'js-cookie';

const TOKEN_KEY = 'zgsmAdminToken';

export const getHashToken = () => {
    const url = new URL(window.location.href);
    const params = new URLSearchParams(url.search);
    const hasToken = params.get('state');

    return hasToken;
};

export const getToken = () => {
    return Cookies.get(TOKEN_KEY);
};

export const setToken = (token: string) => {
    return Cookies.set(TOKEN_KEY, token);
};

export const clearToken = () => {
    // js-cookie 的 remove 受 path/domain 匹配限制，对网关跨父域共享的
    // zgsmAdminToken（Domain=*.sangfor.com, Path=/）会静默失效。
    // 直接写 document.cookie 强制过期，覆盖常见的 path/domain 组合。
    const expires = 'Thu, 01 Jan 1970 00:00:00 GMT';
    const hostParts = window.location.hostname.split('.');
    const parentDomain = hostParts.length >= 2 ? hostParts.slice(-2).join('.') : '';

    document.cookie = `${TOKEN_KEY}=; expires=${expires}; path=/`;
    document.cookie = `${TOKEN_KEY}=; expires=${expires}; path=${window.location.pathname}`;
    if (parentDomain) {
        document.cookie = `${TOKEN_KEY}=; expires=${expires}; path=/; domain=${parentDomain}`;
        document.cookie = `${TOKEN_KEY}=; expires=${expires}; path=/; domain=.${parentDomain}`;
    }

    // 兜底：清理前端以默认 path 写入的 cookie
    Cookies.remove(TOKEN_KEY);
};

export const cleanUrlState = () => {
    const url = new URL(window.location.href);
    const params = new URLSearchParams(url.search);

    if (params.has('state')) {
        params.delete('state');
        const newUrl =
            url.origin +
            url.pathname +
            (params.toString() ? '?' + params.toString() : '') +
            url.hash;
        history.replaceState({}, document.title, newUrl);
    }
};

export const tokenManager = {
    getHashToken,
    getToken,
    setToken,
    clearToken,
    cleanUrlState,
    validateToken: (token: string) => {
        return !!token && token.length > 0;
    },
};
