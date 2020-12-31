export function validateNumber(str) {
    if (str){
        const reg = /^[0-9]*$/
     return reg.test(str)
    } else{
        return true
    }
}