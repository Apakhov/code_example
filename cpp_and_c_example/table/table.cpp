#include <iostream>
#include <functional>
#include <vector>


using namespace std;
template <class T>
class Cell{
public:
    bool real;
    function<T ()> func;
    static function<T ()> getFunc(const Cell<T> &first,const Cell<T> &second);
public:
    Cell():func(nullptr){real = true;};
    Cell<T>& operator=(const Cell<T> &obj);
    Cell<T>& operator=(const T &arg);
    const Cell<T> operator-();
    Cell<T>& operator+=(const Cell<T> &obj);
    Cell<T>& operator-=(const Cell<T> &obj);
    Cell<T>& operator*=(const Cell<T> &obj);
    Cell<T>& operator/=(const Cell<T> &obj);
    const Cell<T> operator*(const Cell<T> &other) const;
    const Cell<T> operator+(const Cell<T> &other) const;
    const Cell<T> operator-(const Cell<T> &other) const;
    const Cell<T> operator/(const Cell<T> &other) const;
    Cell<T>& operator+=(const T &arg);
    Cell<T>& operator-=(const T &arg);
    Cell<T>& operator*=(const T &arg);
    Cell<T>& operator/=(const T &arg);
    const Cell<T> operator*(const T &arg) const;
    const Cell<T> operator+(const T &arg) const;
    const Cell<T> operator-(const T &arg) const;
    const Cell<T> operator/(const T &arg) const;
    template <class U> friend const Cell<U> operator/(U arg, const Cell<U> &obj);
    template <class U> friend const Cell<U> operator-(U arg, const Cell<U> &obj);
    explicit operator T() const{ return func();};
};

template <class T>
Cell<T>& Cell<T>::operator=(const T &arg) {
    func = [arg]()-> T { return arg;};
    return *this;
}

template <class T>
const Cell<T> Cell<T> ::operator-() {
    if(this->real){
        Cell<T> temp(*this);
        function<T ()>  *f1 = &func;
        temp.func = [f1]()-> T { return -(*f1)();};
        temp.real = false;
        return temp;
    }
    Cell<T> temp(*this);
    function<T ()>  f1 = func;
    temp.func = [f1]()-> T { return -(f1)();};
    temp.real = false;
    return temp;
}

template <class T>
Cell<T>& Cell<T>::operator+=(const Cell<T> &obj) {
    const function<T ()>  f1 = func;
    if(obj.real){
        const function<T ()>  *f2 = &obj.func;
        func = [f1,f2]()-> T { return (f1)() + (*f2)();};
        return *this;
    }
    const function<T ()>  f2 = obj.func;
    func = [f1,f2]()-> T { return (f1)() + (f2)();};
    return *this;
}

template <class T>
Cell<T>& Cell<T>::operator-=(const Cell<T> &obj) {
    const function<T ()>  f1 = func;
    if(obj.real){
        const function<T ()>  *f2 = &obj.func;
        func = [f1,f2]()-> T { return (f1)() - (*f2)();};
        return *this;
    }
    const function<T ()>  f2 = obj.func;
    func = [f1,f2]()-> T { return (f1)() - (f2)();};
    return *this;
}

template <class T>
Cell<T>& Cell<T>::operator*=(const Cell<T> &obj) {
    const function<T ()>  f1 = func;
    if(obj.real){
        const function<T ()>  *f2 = &obj.func;
        func = [f1,f2]()-> T { return (f1)() * (*f2)();};
        return *this;
    }
    const function<T ()>  f2 = obj.func;
    func = [f1,f2]()-> T { return (f1)() * (f2)();};
    return *this;
}

template <class T>
Cell<T>& Cell<T>::operator/=(const Cell<T> &obj) {
    const function<T ()>  f1 = func;
    if(obj.real){
        const function<T ()>  *f2 = &obj.func;
        func = [f1,f2]()-> T { return (f1)() / (*f2)();};
        return *this;
    }
    const function<T ()>  f2 = obj.func;
    func = [f1,f2]()-> T { return (f1)() / (f2)();};
    return *this;
}


template <class T>
const Cell<T> Cell<T>::operator+(const Cell<T> &other) const {
    if(this->real && other.real){
        Cell<T> temp;
        const function<T ()>  *f1 = &func;
        const function<T ()>  *f2 = &other.func;
        temp.real = false;
        temp.func = [f1,f2]()-> T { return (*f1)() + (*f2)();};
        return temp;
    }
    if(this->real){
        Cell<T> temp;
        const function<T ()>  *f1 = &func;
        const function<T ()>  f2 = other.func;
        temp.real = false;
        temp.func = [f1,f2]()-> T { return (*f1)() + (f2)();};
        return temp;
    }
    if(other.real){
        Cell<T> temp;
        const function<T ()>  f1 = func;
        const function<T ()>  *f2 = &other.func;
        temp.real = false;
        temp.func = [f1,f2]()-> T { return (f1)() + (*f2)();};
        return temp;
    }
    Cell<T> temp;
    const function<T ()>  f1 = func;
    const function<T ()>  f2 = other.func;
    temp.real = false;
    temp.func = [f1,f2]()-> T { return (f1)() + (f2)();};
    return temp;
}

template <class T>
const Cell<T> Cell<T>::operator-(const Cell<T> &other) const {
    if(this->real && other.real){
        Cell<T> temp;
        const function<T ()>  *f1 = &func;
        const function<T ()>  *f2 = &other.func;
        temp.real = false;
        temp.func = [f1,f2]()-> T { return (*f1)() - (*f2)();};
        return temp;
    }
    if(this->real){
        Cell<T> temp;
        const function<T ()>  *f1 = &func;
        const function<T ()>  f2 = other.func;
        temp.real = false;
        temp.func = [f1,f2]()-> T { return (*f1)() - (f2)();};
        return temp;
    }
    if(other.real){
        Cell<T> temp;
        const function<T ()>  f1 = func;
        const function<T ()>  *f2 = &other.func;
        temp.real = false;
        temp.func = [f1,f2]()-> T { return (f1)() - (*f2)();};
        return temp;
    }
    Cell<T> temp;
    const function<T ()>  f1 = func;
    const function<T ()>  f2 = other.func;
    temp.real = false;
    temp.func = [f1,f2]()-> T { return (f1)() - (f2)();};
    return temp;
}

template <class T>
const Cell<T> Cell<T>::operator*(const Cell<T> &other) const {
    if(this->real && other.real){
        Cell<T> temp;
        const function<T ()>  *f1 = &func;
        const function<T ()>  *f2 = &other.func;
        temp.real = false;
        temp.func = [f1,f2]()-> T { return (*f1)() * (*f2)();};
        return temp;
    }
    if(this->real){
        Cell<T> temp;
        const function<T ()>  *f1 = &func;
        const function<T ()>  f2 = other.func;
        temp.real = false;
        temp.func = [f1,f2]()-> T { return (*f1)() * (f2)();};
        return temp;
    }
    if(other.real){
        Cell<T> temp;
        const function<T ()>  f1 = func;
        const function<T ()>  *f2 = &other.func;
        temp.real = false;
        temp.func = [f1,f2]()-> T { return (f1)() * (*f2)();};
        return temp;
    }
    Cell<T> temp;
    const function<T ()>  f1 = func;
    const function<T ()>  f2 = other.func;
    temp.real = false;
    temp.func = [f1,f2]()-> T { return (f1)() * (f2)();};
    return temp;
}

template <class T>
const Cell<T> Cell<T>::operator/(const Cell<T> &other) const {
    if(this->real && other.real){
        Cell<T> temp;
        const function<T ()>  *f1 = &func;
        const function<T ()>  *f2 = &other.func;
        temp.real = false;
        temp.func = [f1,f2]()-> T { return (*f1)() / (*f2)();};
        return temp;
    }
    if(this->real){
        Cell<T> temp;
        const function<T ()>  *f1 = &func;
        const function<T ()>  f2 = other.func;
        temp.real = false;
        temp.func = [f1,f2]()-> T { return (*f1)() / (f2)();};
        return temp;
    }
    if(other.real){
        Cell<T> temp;
        const function<T ()>  f1 = func;
        const function<T ()>  *f2 = &other.func;
        temp.real = false;
        temp.func = [f1,f2]()-> T { return (f1)() / (*f2)();};
        return temp;
    }
    Cell<T> temp;
    const function<T ()>  f1 = func;
    const function<T ()>  f2 = other.func;
    temp.real = false;
    temp.func = [f1,f2]()-> T { return (f1)() / (f2)();};
    return temp;
}

template <class T>
Cell<T>& Cell<T>::operator+=(const T &arg){
    const function<T ()>  f1 = func;
    func = [f1,arg]()-> T { return (f1)() + arg;};
    return *this;
}

template <class T>
Cell<T>& Cell<T>::operator-=(const T &arg){
    const function<T ()>  f1 = func;
    func = [f1,arg]()-> T { return (f1)() - arg;};
    return *this;
}

template <class T>
Cell<T>& Cell<T>::operator*=(const T &arg){
    const function<T ()>  f1 = func;
    func = [f1,arg]()-> T { return (f1)() * arg;};
    return *this;
}

template <class T>
Cell<T>& Cell<T>::operator/=(const T &arg){
    const function<T ()>  f1 = func;
    func = [f1,arg]()-> T { return (f1)() / arg;};
    return *this;
}

template <class T>
const Cell<T> Cell<T>::operator+(const T &arg) const{
    if(this->real){
        Cell<T> temp;
        const function<T ()>  *f1 = &func;
        temp.real = false;
        temp.func = [f1,arg]()-> T { return (*f1)() + arg;};
        return temp;
    }
    Cell<T> temp;
    const function<T ()>  f1 = func;
    temp.real = false;
    temp.func = [f1,arg]()-> T { return (f1)() + arg;};
    return temp;
}

template <class T>
const Cell<T> Cell<T>::operator-(const T &arg) const{
    if(this->real){
        Cell<T> temp;
        const function<T ()>  *f1 = &func;
        temp.real = false;
        temp.func = [f1,arg]()-> T { return (*f1)() - arg;};
        return temp;
    }
    Cell<T> temp;
    const function<T ()>  f1 = func;
    temp.real = false;
    temp.func = [f1,arg]()-> T { return (f1)() - arg;};
    return temp;
}

template <class T>
const Cell<T> Cell<T>::operator*(const T &arg) const{
    if(this->real){
        Cell<T> temp;
        const function<T ()>  *f1 = &func;
        temp.real = false;
        temp.func = [f1,arg]()-> T { return (*f1)() * arg;};
        return temp;
    }
    Cell<T> temp;
    const function<T ()>  f1 = func;
    temp.real = false;
    temp.func = [f1,arg]()-> T { return (f1)() * arg;};
    return temp;
}

template <class T>
const Cell<T> Cell<T>::operator/(const T &arg) const{
    if(this->real){
        Cell<T> temp;
        const function<T ()>  *f1 = &func;
        temp.real = false;
        temp.func = [f1,arg]()-> T { return (*f1)() / arg;};
        return temp;
    }
    Cell<T> temp;
    const function<T ()>  f1 = func;
    temp.real = false;
    temp.func = [f1,arg]()-> T { return (f1)() / arg;};
    return temp;
}

template <class T>
const Cell<T> operator+(T arg, const Cell<T> &obj){
    return obj + arg;
}

template <class T>
const Cell<T> operator-(T arg, const Cell<T> &obj){
    if(obj.real){
        Cell<T> temp;
        const function<T ()>  *f1 = &obj.func;
        temp.real = false;
        temp.func = [f1,arg]()-> T { return arg - (*f1)();};
        return temp;
    }
    Cell<T> temp;
    const function<T ()>  f1 = obj.func;
    temp.real = false;
    temp.func = [f1,arg]()-> T { return arg - (f1)();};
    return temp;
}

template <class T>
const Cell<T> operator*(T arg, const Cell<T> &obj){
    return obj * arg;
}

template <class T>
const Cell<T> operator/(T arg, const Cell<T> &obj){
    if(obj.real){
        Cell<T> temp;
        const function<T ()>  *f1 = &obj.func;
        temp.real = false;
        temp.func = [f1,arg]()-> T { return arg / (*f1)();};
        return temp;
    }
    Cell<T> temp;
    const function<T ()>  f1 = obj.func;
    temp.real = false;
    temp.func = [f1,arg]()-> T { return arg / (f1)();};
    return temp;
}

template<class T>
Cell<T> &Cell<T>::operator=(const Cell<T> &obj) {
    if(obj.real){
        const function<T ()>  *f2 = &obj.func;
        func = [f2]()-> T { return (*f2)();};
        return *this;
    }else{
        const function<T ()>  f2 = obj.func;
        func = [f2]()-> T { return (f2)();};
        return *this;
    }
}

template <class T>
class SuperCalc {
private:
    int shift;
    vector<Cell<T>> matrix;
public:
    SuperCalc(int m, int n);
    Cell<T>& operator() (int i, int j);
};

template <class T>
SuperCalc<T>::SuperCalc(int m, int n):shift(n), matrix(m*n) {
}

template <class T>
Cell<T>& SuperCalc<T>::operator()(int i, int j) {
    return matrix[i*shift + j];
}



int main()
{
    SuperCalc<int> sc(1, 3);
    sc(0, 2) = -sc(0, 0) * sc(0, 1);
    sc(0, 1) = 300 - sc(0, 0);
    sc(0, 0) = 100;
    cout << (int)sc(0, 2) << endl;
    return 0;
}





