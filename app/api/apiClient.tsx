class API_CONFIG {
    static BASE_URL = __DEV__ 
        ? 'http://localhost:8080/api/v1/'
        : 'https://api.fitup.com/api/v1/';  
    static TIMEOUT = 10000; 
    
}



export { API_CONFIG };