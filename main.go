package main

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"os"
)

type DdnsParam struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Comment string `json:"comment"`
	Proxied *bool  `json:"proxied"`
	Ttl     int    `json:"ttl"`
}

func main() {
	go func() {
		ddns, err := Ddns("newzhxu.com", DdnsParam{
			Type:    "A",
			Name:    "test.newzhxu.com",
			Content: "1.2.3.4",
			Comment: "test",
			Ttl:     1,
			Proxied: cloudflare.BoolPtr(false),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(ddns)
	}()

	select {}
}
func Ddns(zoneName string, param DdnsParam) (cloudflare.DNSRecord, error) {
	api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		fmt.Println(err)
		return cloudflare.DNSRecord{}, err
	}
	zoneId, err := api.ZoneIDByName(zoneName)
	if err != nil {
		fmt.Println(err)
		return cloudflare.DNSRecord{}, err
	}
	ctx := context.Background()
	records, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneId), cloudflare.ListDNSRecordsParams{
		Name: param.Name,
		Type: param.Type,
	})
	if err != nil {
		fmt.Println(err)
		return cloudflare.DNSRecord{}, err
	}
	if len(records) == 0 {
		fmt.Println("create")
		record, err := api.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneId), cloudflare.CreateDNSRecordParams{
			Type:    param.Type,
			Name:    param.Name,
			Content: param.Content,
			Comment: param.Comment,
			Proxied: param.Proxied,
			TTL:     param.Ttl,
		})
		if err != nil {
			fmt.Println(err)
			return cloudflare.DNSRecord{}, err
		}
		return record, nil

	} else {
		fmt.Println("update")
		record, err := api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneId), cloudflare.UpdateDNSRecordParams{
			Type:    param.Type,
			Name:    param.Name,
			Content: param.Content,
			Comment: param.Comment,
			Proxied: param.Proxied,
			TTL:     param.Ttl,
			ID:      records[0].ID,
		})
		if err != nil {
			fmt.Println(err)
			return cloudflare.DNSRecord{}, err
		}
		return record, nil
	}

}
