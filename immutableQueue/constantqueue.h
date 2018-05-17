//
// Created by mikhail on 04.05.18.
//

#ifndef LAB8_CONSTANTQUEUE_H
#define LAB8_CONSTANTQUEUE_H

#include <memory>
#include <stack>

template <typename T>
class ImmutableElem{
private:
    T elem;
    std::shared_ptr<ImmutableElem<T>> next;
public:
    T getElem() {
        return elem;
    }

    ImmutableElem<T>* getNext() {
        return next.get();
    }
public:
    ImmutableElem(T elem, ImmutableElem<T> *nextElem);
};

template <typename T>
class ImmutableStack{
private:
    int size;
public:
    int getSize() const {
        return size;
    }

private:
    std::shared_ptr<ImmutableElem<T>> top;
    std::stack<ImmutableStack<T>*> depStack;
private:
    ImmutableStack(T elem, ImmutableElem<T> *next, int size);
    ImmutableStack(ImmutableElem<T> *next, int size);
public:
    ImmutableStack();
    ImmutableStack(const ImmutableStack &obj);
    ImmutableStack<T>& push(T elem);
    ImmutableStack<T>& pop();
    T peek();
    ImmutableStack<T>&operator=(const ImmutableStack<T> &obj);
    virtual ~ImmutableStack();
};


template  <typename T, int N = -1, bool isEndless = N < 0>
class ImmutableQueue{
private:
    ImmutableStack<T> stackEnq;
    ImmutableStack<T> stackDeq;
    T lastIn;
    std::stack<ImmutableQueue<T,N, isEndless>*> depStack;
private:
    ImmutableQueue(ImmutableStack<T> stackEnq, ImmutableStack<T> stackDeq, T lastIn);


public:
    ImmutableQueue();
    ImmutableQueue<T,N, isEndless>& push(T elem);
    ImmutableQueue<T,N, isEndless>& pop();
    T peek();
    ImmutableQueue<T,N, isEndless>&operator=(const ImmutableQueue<T, N> &obj);
    virtual ~ImmutableQueue();
    int getSize() const {
        return stackEnq.getSize() + stackDeq.getSize();
    }

};

#endif //LAB8_CONSTANTQUEUE_H
