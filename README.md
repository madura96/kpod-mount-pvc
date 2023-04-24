## Introduction 

tooling/kpod-mount-pvc is used to mount an existing volume on a pod


## Specification:  
Create a pod and mount the volume, the script takes these parameters:
1- command takes pvc name as param 
2- create the pod with busybox container (by default) 
3- the image can be changed by giving an argument --image <container-image-of-tools>
4- constitute ref images
5- make a dump of the db
