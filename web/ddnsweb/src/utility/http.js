
import axios from 'axios';
import Qs from 'qs';
import router from './router'
axios.defaults.timeout = 5000;
axios.defaults.withCredentials = true;
axios.defaults.baseURL = process.env.VUE_APP_API_URL;

let GetTocken = (function () {
    return sessionStorage.getItem("__jwt__");
});


//http request 拦截器
axios.interceptors.request.use(
    config => {
        config.headers = {
            'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8',
        }
        return config;
    },
    error => {
        return Promise.reject(error);
    }
);

//http response 拦截器
axios.interceptors.response.use(
    response => {
        return response;
    },
    error => {
        console.log(error.response)
        if (error.response && error.response.status == 403) {
            //router.push("/login");

        }
        return Promise.reject(error)
    }
)


/**
 * 封装get方法
 * @param url
 * @param data
 * @returns {Promise}
 */

export function get(url, params = {}) {

    return new Promise((resolve, reject) => {
        axios.get(url, {
            params: params
        })
            .then(response => {
                if (response.status == 200) {
                    // sendHttpMsg("status 200")

                    resolve(response.data);
                }
            })
            .catch(err => {
                reject(err)
            })
    })
}


/**
 * 封装post请求
 * @param url
 * @param data
 * @returns {Promise}
 */

export function post(url, data = {}, type = {}) {
    data = Qs.stringify(data)
    return new Promise((resolve, reject) => {
        axios.post(url, data, type)
            .then(response => {
                if (response.status == 200) {
                    resolve(response.data);
                }
            }, err => {
                reject(err)
            })
    })
}

/**
 * 封装patch请求
 * @param url
 * @param data
 * @returns {Promise}
 */

export function patch(url, data = {}) {
    data = Qs.stringify(data)
    return new Promise((resolve, reject) => {
        axios.patch(url, data)
            .then(response => {
                if (response.status == 200) {
                    resolve(response.data);
                }
            }, err => {
                reject(err)
            })
    })
}

/**
 * 封装put请求
 * @param url
 * @param data
 * @returns {Promise}
 */

export function put(url, data = {}) {
    data = Qs.stringify(data)
    return new Promise((resolve, reject) => {
        axios.put(url, data)
            .then(response => {
                if (response.status == 200) {
                    resolve(response.data);
                }
            }, err => {
                reject(err)
            })
    })
}


/**
 * 封装delete请求
 * @param url
 * @param data
 * @returns {Promise}
 */

export function del(url, data = {}) {
    // data = Qs.stringify(data)
    return new Promise((resolve, reject) => {
        axios.delete(url, {
            params: data
        })
            .then(response => {
                if (response.status == 200) {
                    resolve(response.data);
                }
            }, err => {
                reject(err)
            })
    })
}
