#include <stddef.h>
#include <linux/bpf.h>
#include <linux/in.h>
#include <linux/if_ether.h>
#include <linux/if_packet.h>
#include <linux/ip.h>
#include <linux/ipv6.h>
#include <linux/icmpv6.h>
#include "bpf_helpers.h"
#include "bpf_endian.h"
/* Defines xdp_stats_map from packet04 */
#include "../common/xdp_stats_kern_user.h"
#include "../common/xdp_stats_kern.h"

struct vlan_hdr {
	__be16	h_vlan_TCI;
	__be16	h_vlan_encapsulated_proto;
};

SEC("xdp/xdp_ip_filter")
int xdp_ip_filter(struct xdp_md *ctx) {
    void *end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;
    __u64 offset;
    __u16 eth_type;

    struct ethhdr *eth = data;
    offset = sizeof(*eth);

    if (data + offset > end) {
    return XDP_ABORTED;
    }
    eth_type = eth->h_proto;

    /* handle VLAN tagged packet 处理 VLAN 标记的数据包*/
       if (eth_type == bpf_htons(ETH_P_8021Q) || eth_type == bpf_htons(ETH_P_8021AD)) {
          struct vlan_hdr *vlan_hdr;

          vlan_hdr = (void *)eth + offset;
          offset += sizeof(*vlan_hdr);
          if ((void *)eth + offset > end)
               return XDP_ABORTED;
          eth_type = vlan_hdr->h_vlan_encapsulated_proto;
    }

    /* let's only handle IPv4 addresses 只处理 IPv4 地址*/
    if (eth_type == bpf_ntohs(ETH_P_IPV6)) {
        return XDP_DROP;
    }

    struct iphdr *iph = data + offset;
    offset += sizeof(struct iphdr);
    /* make sure the bytes you want to read are within the packet's range before reading them
    * 在读取之前，确保你要读取的子节在数据包的长度范围内
    */
    if ((void *)iph + 1 > end) {
        return XDP_ABORTED;
    }

    return xdp_stats_record_action(ctx, XDP_PASS); /* read via xdp_stats */
}

char _license[] SEC("license") = "GPL";