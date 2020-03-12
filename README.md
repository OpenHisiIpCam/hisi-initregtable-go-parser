<!--
<p align="center">
 <img src="images/systemrdl300.png" alt="registers-description-systemdrl">
</p>
-->
<h3 align="center">hisi-initregtable-go-parser</h3>

---

<p align="center">HiSilicon SoC`s U-Boot initial register table parser into human readable format</p>
<p align="center"><em>Part of <a href="https://www.openhisiipcam.org">OpenHisiIpCam</a> project</em></p>

## :eyeglasses: About

This small application has several goals: 
* useful tool for low level HiSilicon ip camera SoCs initialization research (at least it is applicable for chips that use U-Boot 2010.06)
* example how cool and easy we can use machine readable SystemRDL described registers information (more information [here](https://github.com/OpenHisiIpCam/registers-description))

This is example output (left bottom corner is data view before parsing):

![hisi-initregtable-go-parser example screenshot](images/hisi-initregtable-go-parser-example.png)

As you can see, for each item in table of registers initilization we have information about:
* **operation** (WRITE, READ or DELAY)
* **register**, that is used by operation
* register`s **fields**, that are affected
* each field **value meaning**

## :hammer: Usage

### Build

**TODO**

### Run

```shell
foo@bar:~$ hisi-initregtable-go-parser -help
```

```shell
foo@bar:~$ hisi-initregtable-go-parser -file ./regbins/reg_info_hi3519v101.bin -offset 0 -size 4016
```

## :bulb: How does it work

**TODO**

## :exclamation: Futher information

This work were inspired by @kakigate`s [hisi-initregtable-parser](https://github.com/kakigate/hisi-initregtable-parser)

At the moment parser only partially covers hi3516av200 family (hi3519v101 and hi3516av200 chips), 
but this is only a matter of filling the [register database](https://github.com/OpenHisiIpCam/registers-description),
if you are interested your contribution will be appriciated!
