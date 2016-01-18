#!/bin/bash
####version 1.1.4
#rm -fr /tmp/consul/
pkill consul
#rm -fr /etc/consul.d/
#rm -fr /usr/local/bin/consul
#rm -fr /usr/local/bin/consul-template
rm -fr /usr/local/cmha/
#rm -fr /etc/consul-template.d/
rm -fr /etc/haproxy/
rm -fr /etc/keepalived/
rpm -e haproxy-1.5.2-2.el6.x86_64
rpm -e keepalived-1.2.13-5.el6_6.x86_64
pkill consul-template
pkill haproxy
pkill keepalived
