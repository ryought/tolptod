# データ構造

## Succinct bit vector

長さnのbinary vector B[0,n)を考える。
「B[i,j)に1が何個あるか？」というクエリを高速に返したい。
(B[i]からB[j-1]を順にチェックするとO(j-i)~O(n)だが、これをO(1)で実現したい。)

rankをサポートする簡潔ビットベクトルを使えば実現できる。
B[0:i)に含まれる1の個数を数える rank(1, i) という操作がO(1)で出来るデータ構造。

> B[0:i)に含まれる0の個数は rank(0,i) = i-rank(1,i) と計算できる。

スペースを気にしなければ、rank[i]=rank(0,i) i=0...n をO(n)で前もって計算しておけばできる。
> rank配列は長さn+1なのに注意。

> succinctに、記憶容量をn bitに近づけたければ、rank[i]をlog(n)個に1つcacheしておくようにする。
> (簡潔bit vectorの論文を参照)

例えば
```
           0 1 2 3 4 5 6 7    n=8
B       = [0,1,0,1,1,0,1,0]
rank(0) = [0,1,1,2,2,2,3,3,4]
rank(1) = [0,0,1,1,2,3,3,4,4]
               <--->   [i:j)
               i=2   j=5
```

B[i,j)に1が何個あるかは、B[0:j)にある1の個数からB[0:i)にある1の個数を引けば良いので、rank(1,j)-rank(1,i)でO(1)で求められる。

## Wavelet tree/matrix

2012年に提案された
Claude, F., Navarro, G. (2012). The Wavelet Matrix. In: Calderón-Benavides, L., González-Caro, C., Chávez, E., Ziviani, N. (eds) String Processing and Information Retrieval. SPIRE 2012. Lecture Notes in Computer Science, vol 7608. Springer, Berlin, Heidelberg. https://doi.org/10.1007/978-3-642-34109-0_18

先ほどのbinary vectorを、4種類のalphabetの文字列の場合に拡張することを考える。

長さnの文字列S[0,n)を考える。アルファベットは{0,1,2,3}とする。
まず、このアルファベットを、最下位ビットが0か1かで分類する。
```
       !
0 = 0b00
1 = 0b01
2 = 0b10
3 = 0b11
```

```
        0 1 2 3 4 5 6 7
S    = [1 2 0 2 3 1 1 0]
        1 0 0 0 1 1 1 0
        0 1 0 1 1 0 0 0

X[0] = [0 1 2 3 4 5 6 7] ← 添字のリスト
B[0] = [1 0 0 0 1 1 1 0] ← 最下位ビット

       <------> offset[0] = 4
X[1] = [1 2 3 7|0 4 5 6] ← B[0]の順にX[0]を安定ソートした
B[1] = [1 0 1 0|0 1 0 0] ← 第2位ビット

X[2] = [2 7 0 5 6|1 3 4] ← B[1]の順にX[1]を安定ソートした
```


こうすると、

- access
- rank
- intersection

ができる。

### Access
access(i)=S[i]

B[0]上の位置iが、B[1]上でどこに対応するかを知りたい。
つまり、X[0][i]=X[1][i']となるようなi'を求めたい。

B[0][i]=0なら仕切りより左にある。
B[0]の[:i)に0が何個出てくるかは、rank(i)でわかる。
よってi'=rank(B[0], i)

また、B[0][i]=1なら仕切りより右にある。
B[0]の[:i)に1が何個出てくるかは、i-rank(i)でわかる。
B[0]に0が何個出てくるかは、offset[0]として記録してあった。
よってi'=offset[0]+i-rank(B[0], i)

### Rank
rank(c,i)=(文字cがS[:i)に何回出てくるか？)
例えばc=1(bit表現が0b01)の時を考えよう。

「B[d]中の区間[L:R)に対応するS上の要素たちは、先頭d個のbitがcと一致する」という性質を満たすように、B[0]上の[L=0:R=i)から始めて、B[d]まで[L:R)を更新していく。

d-th bitが0ならば
L'= (B[d][0:L)中の0の個数)
R'= (B[d][0:R)中の0の個数)

d-th bitが1ならば
L'= (B[d]中の0の個数) + (B[d][0:L)中の1の個数)
R'= (B[d]中の0の個数) + (B[d][0:R)中の1の個数)

```
            L       R
S    = [1 2 0 2 3 1 1 0]
X[0] = [0 1 2 3 4 5 6 7]
B[0] = [1 0 0 0 1 1 1 0]  d=0

 d-th bit is 0 | 1
               |
X[1] = [1 2 3 7|0 4 5 6]  d=1
B[1] = [1 0 1 0|0 1 0 0]
                  L   R
```

### Top
top(i,j)
区間[i:j)に出現する、一番出現頻度の高い文字を見つける。
Priority Queueを使う

### Intersect
intersect(aL,aR,bL,bR,K)
[aL,aR)と[bL:bR)に共通して登場するK-merがあるか？


## k-mer wavelet matrix

これを拡張して、高々Kのk-merを検索できるデータ構造を作る。
各文字S[i]がS[i:i+k)になっているようなアルファベットがΣ^kとなる文字列を考える。