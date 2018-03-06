# learn_Google_Search_API

## Search APIとは
構造化データを含むDocumentを索引付けするためのモデルを提供する。
Indexを検索し、検索結果を整理して提示することができる。
文字列フィールドの全文検索を行うことができる。
DocumentとIndexは、検索操作に最適化され分散された永続ストアに保存される。
任意の数のDocumentのIndexを作成することができる。

参考 : [Documents and Indexes  |  App Engine standard environment for Go  |  Google Cloud Platform](https://cloud.google.com/appengine/docs/standard/go/search/)

### 全文検索とは
> 全文検索（ぜんぶんけんさく、英: Full text search）とは、コンピュータにおいて、複数の文書（ファイル）から特定の文字列を検索すること。
「ファイル名検索」や「単一ファイル内の文字列検索」と異なり、「複数文書にまたがって、文書に含まれる全文を対象とした検索」という意味で使用される。

引用元 : [全文検索 - Wikipedia](https://ja.wikipedia.org/wiki/%E5%85%A8%E6%96%87%E6%A4%9C%E7%B4%A2)

## データ構造の概念
## Document
一意のIDとユーザーデータを含むfieldのリストを持つObject。
Datastoreで言う所のEntityのようなものだと考えるとわかりやすい。
Documentは、Fieldのリストを含むGo構造体で表される。
ただし、FieldLoadSaverインタフェースを実装する任意のタイプで表現することも可能。
これは、DatastoreにおけるPropertyLoadSaverの話にも似ているかもしれない。(詳しくは[ここ](https://qiita.com/Sekky0905/items/0a4c981ce5dbd0646226)を参照)
FieldLoadSaverを使用すると、DocumentMetadataタイプのドキュメントに対してメタデータを設定できる。
構造体のポインタは、より強く型付けされていて使用しやすいが、FieldLoadSaversはより柔軟に使用することができる。


### Documentの制限
ドキュメントの最大サイズは、1 MBまでなので注意!

### Documentの識別
Index内の全てのDocumentには、一意の識別子(ID)かdocIDが必要。
IDは、searchを使用しないで、IndexからDocumentを取得するのに必要。
docIDはデフォルトでは、Documentが生成された際にSearch APIによって自動的に生み出される。
Documentを生成する際に、自分でdocIDを指定することもできる。

### docIDの制限
表示可能で印刷可能なASCII文字（ASCIIコード33〜126を含む）のみ使用できる。
500文字以下である必要がある。
`！` で始めることはできない。
`__` で始まり、終わることもできない。
 ** docIDをsearchの操作で含めることはできないので、注意! ** 


## Field
Document内の各属性の様なもの。
Datastoreでいう所のPropertyのようなもの。

Documentには、名前、タイプ、およびそのタイプの単一の値を持つFieldが含まれる。
Documentには、同じ名前と同じタイプの複数のFieldを持つことができる。
→ 複数の値を持つフィールドを表現する方法という感じ。
Documentには、同じ名前と異なるFieldタイプの複数のを含めることもできる。

以下の種類が存在する。
> Atom Field - an indivisible character string
  Text Field - a plain text string that can be searched word by word
  HTML Field - a string that contains HTML markup tags, only the text outside the markup tags can be searched
  Number Field - a floating point number
  Time Field - a time.Time value (stored with millisecond precision)
  Geopoint Field - a data object with latitude and longitude coordinates
  
  引用元 : [Documents and Indexes  |  App Engine standard environment for Go  |  Google Cloud Platform](https://cloud.google.com/appengine/docs/standard/go/search/#Go_Documents)



## Index
検索のためにDocumentを保存する。
Documentのグループは、別々のIndexに入れて管理できる。


参考 : [Documents and Indexes  |  App Engine standard environment for Go  |  Google Cloud Platform](https://cloud.google.com/appengine/docs/standard/go/search/#Go_Indexes)

### Indexの制限
IndexのDocumentは制限なし。
使用できるIndexの数も制限なし。
1つのIndexにあるすべてのDocumentの合計サイズはデフォルトで10GBに制限されている。
ただし、Google Cloud PlatformコンソールのApp Engine Searchページからリクエストを送信すると最大で200GBまで
増やすことが可能。

### Index内のDocumentの取得
1つのDocumentを取得する場合はIDで取得する。
範囲でDocumentを取得する場合は、連続するIDで取得する。
Indexを検索して、Query stringとして指定されたFieldとその値について、
指定された基準を満たすDocumentを検索することもできる。

## Query
Queryについては、細かい決まりなどが存在するので、詳細は[公式で](https://cloud.google.com/appengine/docs/standard/go/search/query_strings#Go_Field_search)
確認いただくとして、ざっくりした概要をまとめた。

Indexを検索するには、Query StringでQueryを作成する。
また、場合によっては追加のオプションを含むQueryを作成する。
Query Stringは、1以上のDocumentのFieldの値の条件を指定する。
Indexを検索すると、Queryを満たすFieldを持つIndexのDocumentのみが返却される。

参考 : [Documents and Indexes  |  App Engine standard environment for Go  |  Google Cloud Platform](https://cloud.google.com/appengine/docs/standard/go/search/#Go_Queries)


### Queryの注意事項
atom、text、HTMLのfieldの検索では大文字と小文字が区別されない。
→全部小文字で検索するのがオススメらしい(公式に書いてあった)


### Queryの制限
最大長は2000文字。

### Global searchとField search

#### Global search
Global searchでは、DocumentのFieldに表示される可能性のある値を指定してDocumentを検索する。
Global searchを実行するには、1つ以上のField値を含むQuery Stringを作成する。
検索アルゴリズムは各値のタイプを認識し、その値を含む可能性のあるすべてのDocumentのFieldを検索する。

##### one-value queriesの話
> If the query string is a word ("red") or a quoted string ("\"red rose\""), search retrieves all documents in an index that have:
  
> a text or HTML field that contains that word or quoted string (matching is case insensitive)
  an atom field with a value that matches the word or quoted string (matching is case insensitive)
  If the query string is a number ("3.14159"), search retrieves all documents that have:
  
> a number field with a value equal to the number in the query (a number field with the value 5 will match the query "5" and "5.0")
  a text or HTML field that contains a token that matches the number as it appears in the query (the text field "he took 5 minutes" will match the query "5" but not "5.0")
  an atom field that literally matches the number as it appears in the query
  If the query string is a date in yyyy-mm-dd form, search retrieves all documents that have:
  
> a date field whose value equals that date (leading zeros in the query string are optional, "2012-07-04" and "2012-7-4" are the same date)
  a text or HTML field that contains a token that literally matches the date as it appears in the query
  an atom field that literally matches the date as it appears in the query
  You can prepend the NOT boolean operator (upper case) to a one word query. The result is a list of documents that do not have any fields that match the query value, 
  according to the same rules. So the query "NOT red" will retrieve all documents that don't have any text or HTML fields that contain "red", or any atom fields with the value "red".
  
> You can prepend the NOT boolean operator (upper case) to a one word query. The result is a list of documents that do not have any fields that match the query value,
 according to the same rules. So the query "NOT red" will retrieve all documents that don't have any text or HTML fields that contain "red",
  or any atom fields with the value "red".

引用元 : [Query Strings  |  App Engine standard environment for Go  |  Google Cloud Platform](https://cloud.google.com/appengine/docs/standard/go/search/query_strings#Go_One-value_queries)

意訳

Query Stringが単語（ "red"）または引用符付き文字列（ "\" red rose \ ""）である場合、searchは以下のようなIndexの全Documentを取得する。

* その単語または引用された文字列を含むテキストフィールドまたはHTMLフィールド（一致する文字は大文字と小文字を区別しない）
* 単語または引用符で囲まれた文字列に一致する値を持つatomフィールド（マッチは大文字と小文字を区別しない）

Query Stringが数値（ "3.14159"）の場合、searchは以下を持つすべてのDocumentを取得する。
* Queryの数値と等しい値を持つ数値Field（値5の数値FieldはQuery "5"と "5.0"に一致する）
* Queryに表示される数字に一致するトークンを含むテキストまたはHTMLフィールド（テキストフィールド "he took 5 minutes"はQuery"5.0"ではなく"5"と一致する）
* Query上の数字と文字通り一致するatomフィールド


 Query Stringがyyyy-mm-dd形式の日付である場合、searchは以下を持つすべてのドキュメントを取得する。

* その値がその日付と等しい日付フィールド（Query Stringの先頭のゼロはオプションなので、 "2012-07-04"および "2012-7-4"は同じ日付として扱われる）
* Queryに現れる日付と文字通り一致するトークンを含むテキストまたはHTMLフィールド
* Queryに現れる日付と文字通り一致するアトムフィールド


NOT論理演算子（大文字）を1ワードのクエリに付加す​​ることができる。
結果は、同じ規則に従って、照会値と一致するフィールドを持たない文書のリストになる。
従って、 "NOT red"というQueryは、テキストや "red"を含むHTMLフィールド、または "red"という値を持つ任意のatomフィールドを持たないすべてのDocumentを取得することになる。


##### Multi-value queriesの話
> You can specify multiple values (separated by spaces) in a global search string. 
The white space between words, quoted strings, numbers, and dates is treated as an implicit AND operator. 
The two search strings below are almost the same; they differ in how global search treats atom fields, which is explained below:



引用元 : [Query Strings  |  App Engine standard environment for Go  |  Google Cloud Platform](https://cloud.google.com/appengine/docs/standard/go/search/query_strings#Go_Multi-value_queries)

意訳
Global searchでは、スペースで区切って複数の値を指定することもできる。
単語、引用符で囲まれた文字列、数字、および日付の間の空白は暗黙のAND演算子として扱われる。
以下の2つの検索文字列はほとんど同じ。Global searchがどのようにatomフィールドを扱うかは以下のように異なる。


```go
query = "small red"
query = "small AND red"
```
引用元 : [Query Strings  |  App Engine standard environment for Go  |  Google Cloud Platform](https://cloud.google.com/appengine/docs/standard/go/search/query_strings#Go_Multi-value_queries)




#### Field search
Field searchは、特定のDocument Fieldの値をField名で検索する。
Field searchのQuery Stringは、フィールド名、関係演算子、およびフィールド値を指定する1つ以上の式で書く。

```go
query = "language = go"
```

参考 : [Query Strings  |  App Engine standard environment for Go  |  Google Cloud Platform](https://cloud.google.com/appengine/docs/standard/go/search/query_strings#Go_Field_search)

簡単にDatastoreとの対応を書くと以下。
Kind : Index, Property : Field, Entity : Document

## type SearchOptions

```go
type SearchOptions struct {
    // Limit is the maximum number of documents to return. The zero value
    // indicates no limit.
    Limit int

    // IDsOnly indicates that only document IDs should be returned for the search
    // operation; no document fields are populated.
    IDsOnly bool

    // Sort controls the ordering of search results.
    Sort *SortOptions

    // Fields specifies which document fields to include in the results. If omitted,
    // all document fields are returned. No more than 100 fields may be specified.
    Fields []string

    // Expressions specifies additional computed fields to add to each returned
    // document.
    Expressions []FieldExpression

    // Facets controls what facet information is returned for these search results.
    // If no options are specified, no facet results will be returned.
    Facets []FacetSearchOption

    // Refinements filters the returned documents by requiring them to contain facets
    // with specific values. Refinements are applied in conjunction for facets with
    // different names, and in disjunction otherwise.
    Refinements []Facet

    // Cursor causes the results to commence with the first document after
    // the document associated with the cursor.
    Cursor Cursor

    // Offset specifies the number of documents to skip over before returning results.
    // When specified, Cursor must be nil.
    Offset int

    // CountAccuracy specifies the maximum result count that can be expected to
    // be accurate. If zero, the count accuracy defaults to 20.
    CountAccuracy int
}
```

### Limit int
返すDocumentの最大数を指定する。0の場合は、制限なしになる。

### IDsOnly bool
searchを行った時に、DocumentのIDだけを返すようにするかどうかの真偽値。

### Sort *SortOptions
Search結果の順番をコントロールする。
詳しくは、 `type SortOptions` の項目を参照のこと。

### Fields []string
DocumentのどのFieldsを結果に含むかを指定する。
この項目を除外した場合は、Document内の全てのFieldを返す。

### Expressions []FieldExpression
返却される各Documentに追加処理を行ったFieldsを指定する。
詳しくは、 `FieldExpression` の項目を参照のこと。


### Facets []FacetSearchOption
検索結果に対して返却されるfacet(一面)情報を制御します。
オプションが指定されていない場合、ファセット結果は返されない。

### Refinements []Facet
返却されるDocumentに特定の値を持ったfacetsを含むようにフィルターをかける。
異なる名前のFacetは分離して適用される。
詳しくは、 `type Facet` の項目を参照のこと。

### Cursor Cursor
次のバッチの開始点として使用する。

### Offset int
結果を返す前に幾つのDocumentをスキップするか指定する。
これを指定する場合は、Cursorはnilである必要がある。


### CountAccuracy int
正確であると期待できる最大の結果カウントを指定する。
0の場合、カウント精度はデフォルトで20になる。

## FieldExpression
```go
type FieldExpression struct {
    // Name is the name to use for the computed field.
    Name string

    // Expr is evaluated to provide a custom content snippet for each document.
    // See https://cloud.google.com/appengine/docs/standard/go/search/options for
    // the supported expression syntax.
    Expr string
}
```
### Name
計算フィールドに使用する名前。

### Expr
各Documentにcustom content snippetを提供するために評価される
。

## type SortOptions
結果がどのように計算されて返されるかをより詳細に制御することができる。

```go
type SortOptions struct {
    // Expressions is a slice of expressions representing a multi-dimensional
    // sort.
    Expressions []SortExpression

    // Scorer, when specified, will cause the documents to be scored according to
    // search term frequency.
    Scorer Scorer

    // Limit is the maximum number of objects to score and/or sort. Limit cannot
    // be more than 10,000. The zero value indicates a default limit.
    Limit int
}
```

### FacetSearchOption
```go
type FacetSearchOption interface {
    // contains filtered or unexported methods
}
```

## type Facet 
```go
type Facet struct {
    // Name is the facet name. A valid facet name matches /[A-Za-z][A-Za-z0-9_]*/.
    // A facet name cannot be longer than 500 characters.
    Name string
    // Value is the facet value.
    //
    // When being used in documents (for example, in
    // DocumentMetadata.Facets), the valid types are:
    //  - search.Atom,
    //  - float64.
    //
    // When being used in SearchOptions.Refinements or being returned
    // in FacetResult, the valid types are:
    //  - search.Atom,
    //  - search.Range.
    Value interface{}
}
```
カテゴリー情報をDocumentに追加するために使用する名前と値のペア。

## Faceted Search
> Faceted search gives you the ability to attach categorical information to your documents. A facet is an attribute/value pair. 
For example, the facet named "size" might have values "small", "medium", and "large."
  
> By using facets with search, you can retrieve summary information to help you refine a query and "drill down" into your results in a series of steps.
  
> This is useful for applications like shopping sites, where you intend to offer a set of filters for customers to narrow down the products that they want to see.
  
> The aggregated data for a facet shows you how a facet's values are distributed. For instance, the facet "size" may appear in many of the documents in your result set. 
The aggregated data for that facet might show that the value "small" appeared 100 times, "medium" 300 times, and "large" 250 times. Each facet/value pair represents a subset of documents in the query result. 
A key, called a refinement, is associated with each pair. You can include refinements in a query to retrieve documents that match the query string and that have the facet values corresponding to one or more refinements.
  
> When you perform a search, you can choose which facets to collect and show with the results, or you can enable facet discovery to automatically select the facets that appear most often in your documents.

引用元 : [Faceted Search  |  App Engine standard environment for Go  |  Google Cloud Platform](https://cloud.google.com/appengine/docs/standard/go/search/faceted_search)

意訳
Faceted Searchで、カテゴリ情報をDocumentに付与することができる。
Facetedは属性/値のペアである。例えば、「サイズ」という名前のファセットは、「小」、「中」、「大」の値を持つことがある。

Facetedを検索で使用すると、要約情報を取得して、クエリを絞り込み、一連の手順で結果を「ドリルダウン」するのに役立つ。

これは、(例えば)ショッピングサイトのようなアプリで、顧客が見たい製品を絞り込むための一連のフィルターを提供しようとする場合に便利である。

Facetedの集計データは、Facetedの値がどのように分散されるかを示す。たとえば、Faceted「サイズ」は、結果セット内の多くのDocumentに表示されることがある。
そのFacetedの集計データは、値「小」が100回、「中」が300回、「大」が250回であることを示している可能性がある。
各Faceted/値のペアは、照会結果内のDocumentのサブセットを表す。refinement,と呼ばれるキーは、各ペアに関連付けられている。
Queryにrefinementを含めると、Query Stringに一致し、1つまたは複数のrefinementに対応するFaceted値を持つDocumentを取得することができる。

検索を実行するとき、収集するFacetedと結果を表示するFacetedを選択するか、facet discovery を有効にしてDocumentで最も頻繁に表示されるFacetedを自動的に選択することが可能。


→　Facetsを使用すると、検索結果に、検索で一致したカテゴリの要約を含めることができ、検索を特定のカテゴリに一致するように制限することができる。


