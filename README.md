# data-sweeper-kube
data-sweeper-kube は、主にエッジコンピューティング環境において、マイクロサービス等が生成した不要なファイルを定期的に削除するマイクロサービスです。

# 概要
data-sweeper-kube は、ファイル名や拡張子によって指定された、/var/lib/aion/Data配下(デフォルト設定)のファイルを、一定の期間(interval)ごと、または指定時刻ごとに削除します。  
必要に応じて、外部の API server が data-sweeper-kube を起動します。http://localhost:8080/sweeper にリクエストを送信することで、ターゲットファイルを指定し、削除させることができます。

# 動作環境
data-sweeper-kube は、Kubernetes および AION 上での動作を前提としています。    
以下の環境が必要となります。  
・OS: Linux OS  
・CPU: ARM/AMD/Intel  
・Kubernetes  
・AION  

# 起動方法
Deployment作成前に削除機能の起動方法を設定してください。
設定を変更する場合は`data-sweeper-kube/k8s/data-sweeper-kube.yaml`を開き、`SWEEP_START_TYPE`、`SWEEP_CHECK_INTERVAL`、`SWEEP_CHECK_ALARM`を変更してください。
デフォルトで指定時刻（0時0分0秒）に起動するよう設定してあります。
data-sweeper.yamlファイルを作成し、削除するファイルを指定してください。
具体的なyamlファイルの記述方法は、`sample.yaml`を参照してください。  
yamlファイルの配置場所は、デフォルトでは`/var/lib/aion/default/config`になっています。

|  name                | デフォルト値           | 備考                                       | 
| :------------------: | ------------------- | :---------------------------------------: | 
| SWEEP_START_TYPE     | "alarm"             | 起動方法（"alarm"または"interval"を入力）     | 
| SWEEP_CHECK_INTERVAL | "600000"              | 起動間隔（単位はミリ秒） | 
| SWEEP_CHECK_ALARM    | "00:00:00"          | 起動時刻("HH:mm:ss"の形式で入力)                     |

上記の環境下で、以下のコマンドを入力してDeploymentを作成してください。

1. Docker image の作成
```
$ cd /path/to/data-sweeper-kube
$ make docker-build
```

2. Deploymentの作成
```
$ cd /path/to/data-sweeper-kube
$ kubectl apply -f ./k8s/data-sweeper-kube.yaml
```

3. Deployment作成後、以下のコマンドでPodが正しく生成されていることを確認
```
$ kubectl get pods
```

## I/O
### Input
　　
#### Data Sweeperの実行間隔定義  

下記は、Data Sweeper の実行間隔を指定しています。  
単位は、ミリ秒です。  
yamlファイルは、k8s/data-sweeper-kube.yaml　にあります。  
```
- name: SWEEP_CHECK_INTERVAL
 value: "600000"
```  

実行時にここで指定されたパラメータのミリ秒分は、削除をしません。  
yamlファイルは、yaml/sample.yml　にあります。  
```
sweepTargets:
  - name: 'image'
    fileExtention: [ 'jpg', 'png' ]
    interval: 600000
  - name: 'movie'
    fileExtention: [ 'mp4', 'png' ]
    interval: 600000
  - name: 'text'
    fileExtention: [ 'txt', 'csv', 'json' ]
    interval: 600000
```  

  
#### Data Sweeper 実行対象からの除外
ここで定義したファイル以外は、削除されてしまいますのでご注意ください。  

適切な Database の名前を入れます。  
yamlファイルは、k8s/data-sweeper-kube.yaml　にあります。  
```
- name: MYSQL_DB_NAME
 value: "Hogehoge"
```  

Database で定義された所定のファイルを exclude します。  
main.go で次のように記載されています。  
　
```
func isExitsInDB(filePath string) bool {
```  

#### 外部のAPI serverから指定する場合
json形式でPOSTリクエストを送信してください。
リクエストの例は以下の通りです。
```
{
    "dir_path": "/var/lib/aion/Data",
    "exclude_files": ['202008201552.jpg', '202008111543.mp4', '202008111544.png'],
    "exclude_file_extensions": ['jpg', 'mp4', 'json']
    "is_recursive": true
}
```
### Output
`/var/lib/aion/Data`配下のファイルが削除されます。  

# 参考  

## 各種設定の変更
k8s/data-sweeper.ymlファイルのパラメーターを変更することで、Inputを指定するyamlファイルの配置場所や、削除対象のディレクトリ、intervalを変更することができます。
### ディレクトリの変更
| volumeMounts/volumes | name   | デフォルト値                 | 備考                                   | 
| :------------------: | :----: | ---------------------------- | :------------------------------------: | 
| volumeMounts         | data   | /var/lib/aion/Data           | 削除対象のディレクトリ　(コンテナ上)     | 
| volumeMounts         | config | /var/lib/aion/config         | yamlファイルの配置場所　(コンテナ上) | 
| volumes              | data   | /var/lib/aion/default/Data   | 削除対象のディレクトリ                 | 
| volumes              | config | /var/lib/aion/default/config | yamlファイルの配置場所                 | 

### intervalの変更
| name                 | default | 
| :------------------: | :-----: | 
| SWEEP_CHECK_INTERVAL | 600000ms    | 

## システム図
![system_image](./document/data-sweeper-kube.jpg)
