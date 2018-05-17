#include <iostream>
#include <vector>
#include <memory>
#include "constantqueue.h"



template <typename T>
ImmutableElem<T>::ImmutableElem(T elem, ImmutableElem<T> *nextElem):next(nextElem) {
    this->elem = elem;
}


template <typename T>
ImmutableStack<T>::ImmutableStack():top(nullptr),depStack() {
    size = 0;
    std::cout << "base created" << std::endl;
}

template <typename T>
ImmutableStack<T>::ImmutableStack(const ImmutableStack &obj):top(nullptr),depStack() {
    *this = obj;
}

template <typename T>
ImmutableStack<T>::~ImmutableStack<T>() {
    while (!depStack.empty()){
        delete depStack.top();
        depStack.pop();
    }
   std::cout << "stack deleted" << std::endl;
}

template <typename T>
T ImmutableStack<T>::peek() {
    if (size == 0)
        throw std::runtime_error("no elements in stack(peek)");
    return top.get()->getElem();
}

template <typename T>
ImmutableStack<T>::ImmutableStack(T elem, ImmutableElem<T> *next, int size):top(new ImmutableElem<T>(elem, next)),depStack() {
    this->size = size;
    std::cout << "depended created(push), size:"<< size << std::endl;
}

template <typename T>
ImmutableStack<T>::ImmutableStack(ImmutableElem<T> *next, int size):top(next->getNext()),depStack() {
    this->size = size;
    std::cout << "depended created(pop), size:"<< size << std::endl;
}

template <typename T>
ImmutableStack<T>& ImmutableStack<T>::push(T elem) {
    depStack.push((new ImmutableStack<T>(elem, this->top.get(),this->size+1)));
    return *depStack.top();
}

template <typename T>
ImmutableStack<T>& ImmutableStack<T>::pop() {
    if (size == 0)
        throw std::runtime_error("no elements in stack(pop)");
    depStack.push((new ImmutableStack<T>(this->top.get(),this->size-1)));
    return *depStack.top();
}

template <typename T>
ImmutableStack<T>& ImmutableStack<T>::operator=(const ImmutableStack<T> &obj) {
    std::cout << "coping" << std::endl;
    if(this->top.get() == obj.top.get()){
        return *this;
    }

    this->size = obj.size;
    this->top = obj.top;

}

template <typename T, int N, bool isEndless>
ImmutableQueue<T, N, isEndless>::ImmutableQueue():stackEnq(),stackDeq(),depStack() {
}

template <typename T, int N, bool isEndless>
ImmutableQueue<T, N, isEndless>::ImmutableQueue(ImmutableStack<T> stackEnq, ImmutableStack<T> stackDeq, T lastIn):stackEnq(stackEnq),stackDeq(stackDeq),depStack(),lastIn(lastIn) {}

template <typename T, int N, bool isEndless>
T ImmutableQueue<T, N, isEndless>::peek() {
    if(!isEndless && stackDeq.getSize() + stackEnq.getSize() == 0){
        throw std::runtime_error("no elements in queue(peek)");
    }

    if (stackDeq.getSize() == 0){
        return lastIn;
    }

    return stackDeq.peek();
}

template <typename T, int N, bool isEndless>
ImmutableQueue<T,N, isEndless>& ImmutableQueue<T, N, isEndless>::push(T elem) {
    if(stackDeq.getSize() + stackEnq.getSize() == N){
        throw std::runtime_error("queue overflow");
    }
    if(stackEnq.getSize() == 0)
        depStack.push((new ImmutableQueue(stackEnq.push(elem), stackDeq, elem)));
    else
        depStack.push((new ImmutableQueue(stackEnq.push(elem), stackDeq, lastIn)));

    return *depStack.top();
}

template <typename T, int N, bool isEndless>
ImmutableQueue<T,N, isEndless>& ImmutableQueue<T, N, isEndless>::pop() {
    std::cout << "sizes: "<< stackEnq.getSize() << stackDeq.getSize() <<std::endl;
    if(stackDeq.getSize() + stackEnq.getSize() == 0){
        throw std::runtime_error("no elements in queue(pop)");
    }

    if(stackDeq.getSize() != 0){
        depStack.push((new ImmutableQueue(stackEnq, stackDeq.pop(),lastIn)));
    } else{
        depStack.push(new ImmutableQueue(this->stackEnq, this->stackDeq, lastIn));
        while (depStack.top()->stackEnq.getSize() != 1){
            depStack.top()->stackDeq = depStack.top()->stackDeq.push(depStack.top()->stackEnq.peek());
            depStack.top()->stackEnq = depStack.top()->stackEnq.pop();
        }
        depStack.top()->stackEnq = depStack.top()->stackEnq.pop();
    }
    return *depStack.top();
}

template <typename T, int N, bool isEndless>
ImmutableQueue<T,N, isEndless>& ImmutableQueue<T, N, isEndless>::operator=(const ImmutableQueue<T, N> &obj) {
    this->stackEnq = obj.stackEnq;
    this->stackDeq = obj.stackDeq;
    this->lastIn = obj.lastIn;
}

template <typename T, int N, bool isEndless>
ImmutableQueue<T,N, isEndless>::~ImmutableQueue<T,N, isEndless>() {
    while (!depStack.empty()){
        delete depStack.top();
        depStack.pop();
    }

    std::cout << "deleted!"<<std::endl;
}

int main() {

    ImmutableStack<int> stack1;
    stack1.push(12);
    stack1 = stack1.push(13).push(15).push(19).pop().pop();
    stack1.push(14);
    std::cout << stack1.peek() << std::endl;



    ImmutableQueue<int, 5> queue1;

    queue1 = queue1.push(-1);
    std::cout << queue1.peek()<<std::endl;


    for(int i = 0; i < 100; i++){
        queue1 = queue1.push(i);
        queue1 = queue1.push(++i);
        std::cout << queue1.peek()<<std::endl;
        queue1 = queue1.pop();
        std::cout << queue1.peek()<<std::endl;

        queue1 = queue1.pop();
    }

    std::cout << queue1.peek()<<std::endl;
    queue1 = queue1.pop();
    ImmutableQueue<long> queue2;
    for(long i = 0; i < 1000; i+= 10){
        queue2 = queue2.push(i);
    }
    for(long i = 0; i < 1000; i+= 10){
        std::cout << queue2.peek()<<std::endl;
        queue2 = queue2.pop();
    }
    queue2 = queue2.push(12);
    queue2 = queue2.push(13);
    queue2 = queue2.pop();
    queue2 = queue2.pop();

    std::cout << "ok"<<std::endl;
    return 0;
}