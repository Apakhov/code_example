#include <iostream>
#include <fstream>
#include <vector>
#include <set>

using namespace std;
class Bigram{
private:
    string word;
public:
    const string &getWord() const {
        return word;
    }

private:
    vector<string> bigs;
public:
    Bigram(string s);
    vector<string> &getBigs(){
        return bigs;
    }
    int getSame(const Bigram &b) const;
    int getAll(const Bigram &b) const;
};

Bigram::Bigram(string s):bigs() {
    set<string> used;
    word = s;
    if(s.length() == 1){
        bigs.emplace_back(s);
        return;
    }
    for (int i = 0; i < s.length() - 1; i++){
        if(used.find(s.substr(i, 2)) == used.end()){
            bigs.emplace_back(s.substr(i, 2));
            used.insert(s.substr(i, 2));
        }
    }
}

int Bigram::getSame(const Bigram &b) const {
    set<string> same;
    for(auto &i : this->bigs){
        same.insert(i);
    }
    int res = 0;
    for(auto &i : b.bigs){
        if(same.find(i) != same.end()){res++;}

    }
    return res;
}

int Bigram::getAll(const Bigram &b) const {
    set<string> same;
    for(auto &i : this->bigs){
        same.insert(i);
    }
    for(auto &i : b.bigs){
        same.insert(i);
    }
    return (int)same.size();
}

int main() {
    ifstream fin("./count_big.txt");
    if (!fin.is_open()){
        cout << "can not open file!\n";
        return 1;
    }
    string buff;
    int k;
    vector<Bigram> words;
    vector<int> freqs;
    while(!fin.eof()){
        fin >> buff;
        fin >> k;
        words.emplace_back(Bigram(buff));
        freqs.emplace_back(k);
    }

    string curWord;
    while (getline(cin, curWord)){
        Bigram curBigram(curWord);
        set<string> curBigs;
        for (auto &i : curBigram.getBigs()){
            curBigs.insert(i);
        }
        int same = 0;
        int all = 100000;
        int freq = 0;
        int startsize = (int)curBigram.getBigs().size();
        string bestRes;
        for(int i = 0; i < words.size(); i++){
            int curSame = 0;
            int curAll = startsize;
            for(auto &j : words[i].getBigs()){
                if(curBigs.find(j) != curBigs.end()){
                    curSame++;
                }else{
                    curAll++;
                }
            }
            int curFreq = freqs[i];
            if(same == 0){
                if(curSame != 0 || curAll < all || (curAll == all && curFreq > freq)){
                    same = curSame;
                    all = curAll;
                    freq = curFreq;
                    bestRes = words[i].getWord();
                }
            } else{
                if(same * curAll < curSame * all || (same * curAll == curSame * all && curFreq > freq)){
                    same = curSame;
                    all = curAll;
                    freq = curFreq;
                    bestRes = words[i].getWord();
                }
            }
        }
        cout << bestRes << endl;
    }

    fin.close(); 
    return 0;
}