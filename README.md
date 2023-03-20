# go-dparm

Control ATA/NVMe and TCG protocol

C++ Version: [jcu-dparm](https://github.com/jc-lab/jcu-dparm)

# WIP

This project is working in progress.

# Supported platforms

## Windows

* ATA passthrough command
* TCG protocol support
    - Support on ATA, SCSI, Windows NVMe driver
    - nvmewin driver is not tested
* If the device ― likely USB Flash Memory ― does not support identify, used STORAGE_DEVICE_DESCRIPTOR Instead.

* NVMe passthrough command is not supported yet

## Linux

* All feature support (sg, nvme driver)
* If the device ― likely USB Flash Memory ― does not support identify, used INQUIRY command Instead.


# example

TODO

```go
...
```


# Notes

* 유용한 정보들 : https://jsty.tistory.com/237?category=881462

* TCG 관련 요약 : https://jsty.tistory.com/239

* TCG Locking 확인 : https://jsty.tistory.com/238