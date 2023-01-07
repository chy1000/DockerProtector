# Introduction

`docker`映射端口会自动在`iptables`建立规则绕过`firewalld`，这样端口将暴露在外网，这对服务器的安全造成了很大的隐患。

本程序可根据用户配置的访问规则，生成`docker`的`iptables`规则只允许特定的`IP`访问`docker`映射端口，解决`firewalld`被绕过造成的安全隐患。
