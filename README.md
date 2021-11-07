# data-sweeper-kube
data-sweeper-kube は、主にエッジコンピューティング環境において、マイクロサービス等が生成した不要なファイルを定期的に削除するマイクロサービスです。

# 概要
data-sweeper-kube は、ファイル名や拡張子によって指定された、/var/lib/aion/Data配下(デフォルト設定)のファイルを、一定の期間(interval)ごと、または指定時刻ごとに削除します。  
必要に応じて、外部の API server が data-sweeper-kube を起動します。http://localhost:8080/sweeper にリクエストを送信することで、ターゲットファイルを指定し、削除させることができます。

# 動作環境
data-sweeper-kube を動作させるためには、以下の環境が必要となります。  

* OS: Linux OS  
* CPU: ARM/AMD/Intel  
* Kubernetes  


# 起動方法
Deployment作成前に削除機能の起動方法を設定してください。
設定を変更する場合は`data-sweeper-kube/yaml/sample.yml`を開き、`SWEEP_START_TYPE`、`SWEEP_CHECK_INTERVAL`、`SWEEP_CHECK_ALARM`を変更してください。
デフォルトで指定時刻（0時0分0秒）に起動するよう設定してあります。
sample.yamlファイルを編集し、削除するファイルを指定してください。
具体的なyamlファイルの記述方法は、`sample.yaml`を参照してください。  
yamlファイルの配置場所は、デフォルトでは`/var/lib/aion/default/config`になっています。

|  name                | デフォルト値           | 備考                                       | 
| :------------------: | ------------------- | :---------------------------------------: | 
| SWEEP_START_TYPE     | "alarm"             | 起動方法（"alarm"または"interval"を入力）     | 
| SWEEP_CHECK_INTERVAL | "600000"              | 起動間隔（単位はミリ秒） | 
| SWEEP_CHECK_ALARM    | "00:00:00"          | 起動時刻("HH:mm:ss"の形式で入力)                     |


# data-sweeper-kube のデプロイ・稼働
data-sweeper-kube の デプロイ・稼働 を行うためには、aion-service-definitions の services.yml に設定する必要があります。

ymlファイル(services.yml)の中身
```
  data-sweeper:
    scale: 1
    startup: yes
    always: yes
    network: NodePort
    ports:
      - name: data-sweeper
        protocol: TCP
        port: xxxx
        nodePort: xxxx
    env:
      TZ: Asia/Tokyo
      MYSQL_USER: "xxxx"
      MYSQL_PASSWORD: "xxxx"
      MYSQL_SERVICE_HOST: "mysql"
      MYSQL_SERVICE_PORT: "xxxx"
      MYSQL_DB_NAME: "xxxx"
    volumeMountPathList:
      # Aion上でVolumeをマウントする場合、k8s上でボリュームを指定する場合と異なり、volumeMountsのパス:volumeのパスという書き方をします。
      # volumeには、yml が配置されている場所を指定します。デフォルトでの配置場所は/var/lib/aion/default/configになっています。
      - /var/lib/aion/config:/var/lib/aion/default/config
```

Deployment作成後、以下のコマンドでPodが正しく生成されていることを確認してください。
```
$ kubectl get pods
```

## I/O
### Input
　　
#### Data Sweeper の実行間隔定義  

下記は、Data Sweeper の実行間隔を指定しています。  
単位は、ミリ秒です。  
yamlファイルは、yaml/sample.yml　にあります。  
```
sweepSettings:
     sweepCheckInterval: 600000
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

適切な Database(MySQL) の名前を入れます。  
yamlファイルは、services.yaml　にあります。  
```
      MYSQL_USER: "xxxx"
      MYSQL_PASSWORD: "xxxx"
      MYSQL_SERVICE_HOST: "mysql"
      MYSQL_SERVICE_PORT: "xxxx"
      MYSQL_DB_NAME: "xxxx"
```  

Database で定義された所定のファイルを exclude します。  
main.go で次のように記載されています。  
　
```
func isExitsInDB(filePath string) bool {
```  

また、接続するMYSQLにdata_sweeper_ignore_tables というテーブルがある場合、そこに記載されているカラムとひもづくファイルを除外することができます。詳しくは[こちら](https://bitbucket.org/latonaio/data-sweeper-kube-sql/src/master/)を参照してください。

#### 外部の API Server から指定する場合
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
sample.ymlファイルのパラメーターを変更することで、Inputを指定するyamlファイルの配置場所や、intervalを変更することができます。
### ディレクトリの変更
| volumeMounts/volumes | name   | デフォルト値                 | 備考                                   | 
| :------------------: | :----: | ---------------------------- | :------------------------------------: | 
| volumeMounts         | config | /var/lib/aion/config         | yamlファイルの配置場所　(コンテナ上) | 
| volumes              | config | /var/lib/aion/default/config | yamlファイルの配置場所                 | 

### intervalの変更
| name                 | default | 
| :------------------: | :-----: | 
| SWEEP_CHECK_INTERVAL | 600000ms    | 

## Output
/var/lib/aion/Data配下のファイルが削除されます。

# システム図
![system_image](./document/data-sweeper-kube.jpg)
