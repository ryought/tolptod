# BWT

```py
s = 'mississippi#'
for (p, i) in sorted([(s[i:] + s[:i], i) for i in range(len(s))]):
    print(p, i)
```

```
x  S[i:]    BWT
x  F          L i  LCP
0  #mississippi 11 0
1  i#mississipp 10 0
2  ippi#mississ 7  1
3  issippi#miss 4  1
4  ississippi#m 1  4
5  mississippi# 0  0
6  pi#mississip 9  0
7  ppi#mississi 8  1
8  sippi#missis 6  0
9  sissippi#mis 3  2
10 ssippi#missi 5  1
11 ssissippi#mi 2  3
```

```
0  mississippi#
1  ississippi#m
2  ssissippi#mi
3  sissippi#mis
4  issippi#miss
5  ssippi#missi
6  sippi#missis
7  ippi#mississ
8  ppi#mississi
9  pi#mississip
10 i#mississipp
11 #mississippi
```

「V < W なら、aV < aW」
S[i:]がpi番目で、L[pi]がk番目のaだとする
→S[i-1:]はaで始まるk番目のsuffix→SA中である領域内にある。

LF-mapping: LF(i)=j s.t. SA[j]=SA[i]-1
つまり、S[k:]の順位iから、S[k-1:]の順位jを返す。

C[c] = cより小さい文字が何個出てきているか？ (SA中で、cで始まるsuffixの最初の位置)
LF(i)=C[L[i]]+rank(L[i], L, i)
L[:i]にL[i]が何回出てきているか？が分かれば良い。
S.query(T)がO(|T|)でできる