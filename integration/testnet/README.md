### To launch a local privatenet with 4 validators ###

First follow the steps [here](https://docs.dnerochain.org/docs/setup) to compile the latest code of the `privatenet` branch. Next, run the commands below to launch the privatenet with 4 validators:

```
cd $DNERO_HOME
cp -r ./integration/testnet ../testnet
mkdir ../testnet/logs

# In terminal 1
dnero start --config=../testnet/node1 |& tee ../testnet/logs/node1.log

# In terminal 2
dnero start --config=../testnet/node2 |& tee ../testnet/logs/node2.log

# In terminal 3
dnero start --config=../testnet/node3 |& tee ../testnet/logs/node3.log

# In terminal 4
dnero start --config=../testnet/node4 |& tee ../testnet/logs/node4.log
```
