# tolptod

dotplotを高速にinteractiveに描くツール

## 問題設定

- 長さnの配列S,T
- そのうち、長さmの区間S[x:x+m]とT[y:y+m]のdotplotを描く
- 画面にはs (bp/px)のスケールで描画される、つまり画面サイズw=m/sとする。

## 一番簡単な方法

- suffix array SA(S)を作る。(SAISでO(n))
- T[y:y+m]に含まれる各k-mer q=T[y+i:y+i+k]に対して、SA.LookUp(q)してS[x:x+m]のどこに当たるかを調べる。(m個のk-merそれぞれに対して、O(log(n))かけて二分探索する)

前処理O(n)、各描画あたりO(mlog(n))で色を塗るべきマス目がわかる。

あるいは、
- S#T$のsuffix array SAとその逆関数ISAを作る
- T[y:y+m]に対して、ISAでSA上のrankを調べる。それと同じprefixを持つsuffixをSA上で見つける。

mが小さく、sが1に近い場合はこれが一番早そうだ。
しかし、m/sが大きい場合(例えば1Gbpの領域を1000pxの画面に表示するような時)は、前計算しておくとかでもうちょっと高速にできそうな気がする。
描画する際に最低でもO(w^2)はかかるが、各pixelあたり定数(少なくともO(m)やO(n)以下)に抑えたい。

## 他の手法

### 1: k-mer matchやMEMのリストをもっておく

k-mer matchやmaximal exact match (MEM)の個数(種類数)は高々O(n)個(S#T$としてconcatしたsuffix treeのnodeのsubsetなので)
k-mer数は大体O(n)だろうけど、match (i.e. repetitive) k-mer数はそれよりは少ない。
kは1〜100ぐらいと仮定して良いから、各kについて前処理をするのは大丈夫。
例えばk=16で固定されていることを考える。

match kmerごとの出現位置のソート済みリスト、位置ごとにkmerのリストを持っておく
```
X = {
    ACTT: [10, 20, 50, 100],
    GGAA: [11, 80],
    ...
}
Y = [
    ACTT or (pos, k-mer id)
    CTTG
    TTGA
    ...
]
```

- mの区間に入っているmatch k-merの一覧を取得: Yの二分探索
- その各match k-merに対して、S[x:x+m]とT[y:y+m]に入っているかを二分探索で調べれば、binを色塗るべきかどうかがわかる。


区間内のmatch k-mer数に線形になる。match数が少ない場合は高速になるが、全対全のdotplotの時はやはり低速になってしまう。(n個ぐらいのmatch k-merが生じてしまう。)

### 2: binごとの結果をcacheしておく

s=1000の時、s=1の情報は考えなくてもよい？
(例えば世界地図を表示する時、交差点の形は意識しなくて良い。あるいはHi-Cマトリクスのように、引いてみた時のmatrixをprecomputeしておく)
binに含まれる


## 2つの区間の間にmatchがあるかどうかをO(1)で返す

例えば長い文字列上の大きい2区間に、同じ文字が出現しているか？は、定数時間で返せそう。(A/C/G/Tが出ているかだけ分かればよくて、元の文字列をチェックする必要はない。) k-merも、kが短い場合は定数(O(k)ぐらい)で返せそうじゃないか？と思う。

S#T$上の位置にmatch k-mer idが割り振られている。(無い位置は-1とする)
ある2区間 [i:j] [i':j'] (つまりS[x:x+m]とT[y:y+m])に、同じidが存在するか？を返したい。
wavelet treeで、log(u) u:アルファベットサイズ ぐらいで返せる

> New algorithms on wavelet trees and applications to information retrieval が詳しい。

wavelet matrixを応用したk-merのDBを作った。wavelet.mdを参照。