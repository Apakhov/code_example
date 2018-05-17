#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#define MAX_LEN 1000000

int fenwick_tree_quiry(int *tree, int left, int right){
        int res = 0;
        for(; right >= 0; right = (right & (right + 1)) - 1)
		res ^= tree[right];
	for(left-- ; left >= 0; left = (left & (left + 1)) - 1)
		res ^= tree[left];
	return res;
}

void fenwick_tree_upd(int *tree, int nel, int pos, int delta){
	for(; pos < nel; pos = (pos | (pos + 1)))
		tree[pos] ^= delta;
}

int rec_fenwick_build(int *arr, int nel, int left, int right, int* tree){
	int sum = 0;
	int bound = right < nel ? right : nel;
	while(left < bound){
		int mid = (left + right) / 2;
		sum ^= rec_fenwick_build(arr, nel, left, mid, tree);
		left = mid + 1;
	}
	if (right < nel){
		sum ^= arr[right];
		tree[right] = sum;
	}
	return sum;
} 

int *fenwick_tree_build(int *arr, int nel){
	int *tree = malloc(nel * sizeof(int));
        if (tree == NULL){
                printf("malloc error");
                return NULL;
        }
	int right  = 1;
	while (right < nel)
		right *= 2;
	rec_fenwick_build(arr, nel, 0, right - 1, tree);
	return tree;
}

int main() {
	int *quasistr = malloc(MAX_LEN * sizeof(int));
        if (quasistr == NULL){
                printf("malloc error");
                return 0;
        }
	char new_char;
	int len = 0;
	while((new_char = getchar()) != '\n') {
		quasistr[len++] = (1 << (new_char - 'a'));
	}
        
	int i = 0;
	int *tree = fenwick_tree_build(quasistr, len);
	if (tree == NULL){
                return 0;       
	}
	int amt_com;
	scanf("%d ", &amt_com);
	char com[3];
	int left;
	int right;
	int pos;
	int res;
	for(int i = 0; i < amt_com; i++){
		scanf("%s ", com);
		if (com[0] == 'H'){
			scanf(" %d %d ", &left, &right);
		res = fenwick_tree_quiry(tree, left, right);
		if (res & (res - 1))
			printf("NO\n");
		else
			printf("YES\n");
		} else {
			scanf(" %d ", &pos);
			while((new_char = getchar()) != '\n') {
				fenwick_tree_upd(tree, len, pos, (1 << (new_char - 'a')) ^ (quasistr[pos]));
				quasistr[pos++] = (1 << (new_char - 'a'));
			}
		scanf(" ");
		}
	}
	free(quasistr);
	free(tree);
        
        
        
	return 0;
}
