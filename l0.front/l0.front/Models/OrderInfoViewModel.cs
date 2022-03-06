using System;
using System.Collections.Generic;

namespace l0.front.Models
{
    public class OrderInfo
    {
        public string order_uid { get; set; }
        public string track_number { get; set; }
        public string entry { get; set; }
        public Delivery delivery { get; set; }
        public Payment payment { get; set; }
        public IList<Item> items { get; set; }
        public string locale { get; set; }
        public string internal_signature { get; set; }
        public string customer_id { get; set; }
        public string delivery_service { get; set; }
        public string shardkey { get; set; }
        public int sm_id { get; set; }
        public DateTime date_created { get; set; }
        public string oof_shard { get; set; }

        public struct Delivery
        {
            public string name { get; set; }
            public string phone { get; set; }
            public string zip { get; set; }
            public string city { get; set; }
            public string address { get; set; }
            public string region { get; set; }
            public string email { get; set; }
        }

        public class Payment
        {
            public string transaction { get; set; }
            public string request_id { get; set; }
            public string currency { get; set; }
            public string provider { get; set; }
            public int amount { get; set; }
            public int payment_dt { get; set; }
            public string bank { get; set; }
            public int delivery_cost { get; set; }
            public int goods_total { get; set; }
            public int custom_fee { get; set; }
        }

        public class Item
        {
            public int chrt_id { get; set; }
            public string track_number { get; set; }
            public int price { get; set; }
            public string rid { get; set; }
            public string name { get; set; }
            public int sale { get; set; }
            public string size { get; set; }
            public int total_price { get; set; }
            public int nm_id { get; set; }
            public string brand { get; set; }
            public int status { get; set; }
        }
    }
}