

class API_CONFIG {
    static BASE_URL = __DEV__ 
        ? 'http://localhost:8080'    // Development - API Gateway
        : 'https://api.fitup.com';   // Production
    static TIMEOUT = 10000; 
}

export { API_CONFIG };