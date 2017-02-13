0. assigner 简介
	.assigner是一款基于etcd和pipework的ip自动分配工具。可以从应用的ip池中自动获取ip并分配给pod或者container

1. assigner 安装
    	.*以下步骤每个物理机均需要实施
    	.安装/更新 iproute 到指定版本
    	.将assigner程序、pipework程序放置到/usr/local/bin或者其他能够直接命令启动的路径下
   	 .使用kubelet-yxb代替kubelet。*kubelet-yxb在ETCDPATH未设置时与原生kubelet功能完全一致，即修改是非侵入性的。但是一旦设置ETCDPATH，那将是侵入性的，必须使用assigner
    	.添加一个环境变量，名字必须为 ETCDPATH

2. assigner 使用指南
    	.manage，在一个etcd集群中只需要运行一次，将桥接和子网掩码上传到etcd
    	.import，每一个app仅需要运行一次，填写相关yaml文件即可。如果要更新信息（特别是存在ip不可用的情况），请使用etcdctl手动删除该app目录后重新import。
    	.delete，用于删除ip和container的对应关系，可以使用ip或者containerid删除对应ip。仅支持解锁ip，不可删除其他键值对。已继承在kubelet-yxb程序中，使用kube正常启停不需要手动执行。如果非正常启停则需要手动删除。
	.get，已集成在kubelet-yxb程序中，一般不需要手动执行。
    	.query，用于根据containerid查询ip
    	.help，用于查询assigner命令格式

3. 使用kubelet-yxb调用assigner
	.使用kubelet-yxb调用assigner有两个必要条件，其一是设置了 ETCDPATH，其二是检测到pod的env内有key为APP的键值对。这俩条件缺一都不会调用assigner，而是正常使用kubelet原有功能。
	.定义kubernetes的rc（pod）文件时，添加一个env，其key必须为 APP ，其value为对应的app的信息，且必须为 app:component 格式。如果已经检测到APP，且APP设置不正常，可能会存在pod没有任何ip的问题，即分配ip失败。

4. 使用assigner前提条件
	.所有物理机均已设置桥接，且以桥接为docker的桥而非docker0，详见http://blog.163.com/hk_bs/blog/static/245038011201631931849968/

5. assigner详细设计
    	.manage, 仅与etcd交互，manage命令调用etcdclient的go语言包，在etcd下创建目录/assigner/ips/，/assigner/apps/ ,/assigner/ids/ 和一个键值对 /assigner/config，内容为 桥接;子网掩码。桥接是指传参给docker -b的桥的名字，子网掩码是32以内数字，比如255.255.255.0是24
    	.import, 仅与etcd交互，import命令调用go语言yaml包解析读取到的app的yaml文件，并在/assigner/apps/目录下创建/appname/componetname/目录，该目录下为对应的可分配ip信息的键值对，key为ip，内容为网关
    	.delete, 仅与etcd交互，delete命令输入参数可以是ip，也可以是containerid。当输入参数是ip时，直接找到/assigner/ips下对应ip并删除，达到解锁目的。删除ip后，根据value删除对应的/assigner/ids/下的id。当输入是containerid时，先查询/assigner/ids下索引是否有该id（必须完全匹配），有的话直接使用id下存储的ip即可。没有的话对/assigner/ips下的所有value，用正则表达式{containerid}*来匹配，删除所有匹配到的ip。这样匹配有利的地方在于，kubernetes一直是用container的全长度id进行操作，所以存在etcd里的key很长，这时使用标准长度（docker ps）进行删除或者查询是肯定可以查到的。但是不利的地方在于，人为故意使用过短的containerid可能会引发异常。这里人为操作包括手动使用get命令并使用正常长度id（会导致kubernetes删除不掉，因为kube使用全长度，会匹配不到）；人为进行删除操作，并使用极短的id，匹配到多个ip，都会删除。
    	.query, 仅与etcd交互，输入可以是ip和containerid，设计见delete，不同仅在于不删除而返回查询值。(使用ip查ip很奇怪倒是,相当于是用来确认ip是否存在)   
    	.get,*get会调用pipework。get命令使assigner的核心，该命令在获取到输入参数后，解析app:componet后去etcd对应路径下对所有ip进行一次轮询测试（锁测试）。该测试为：随机一个初始值，以初始值为下标拿到对应的ip（etcd目录下可以认为是个数组，每次内容一致，随机数可以避免每次都从同一个ip开始测试），尝试在/assigner/ips目录下创建key，如果key存在则创建失败继续下一个值，如果创建成功，value为对应的containerid，说明该值为contianer对应的ip。获取到ip后，用pipework以及对应参数为container设置veth，设置完过5s后ping该ip，如果ping通则分配成功，在/assigner/ids下建立key为id,value为ip的键值对(用于做索引减少轮询次数), 然后程序退出。如果ping失败，并且尝试上述pipework和ping命令3次失败，则分配失败，删除/assigner/ips下对应的key（veth并未删除，是个问题）后，程序退出。

6. kubelet改动详解
    	.获取ip，https://github.com/yansmallb/kubelet-yxb/blob/kubelet-yxb/README.
    	.删除ip，直接修改在killContainer的末尾，在killContainer函数退出前，判断ETCDPATH是否设置，如果设置，那么就调用assigner delete containerid etcdpash的方式，释放锁。该调用如果containerid存在，删除key是必须的，当containerid不存在时由于本身key就不存在不会引起任何异常，所以没有其他任何异常处理。
