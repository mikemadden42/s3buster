package main

import "testing"

func TestParseListBucketResult(t *testing.T) {
	const body = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Name>example-bucket</Name>
  <Prefix></Prefix>
  <Marker></Marker>
  <MaxKeys>1000</MaxKeys>
  <IsTruncated>false</IsTruncated>
  <Contents>
    <Key>photo.jpg</Key>
    <LastModified>2024-01-01T00:00:00.000Z</LastModified>
    <ETag>"abc123"</ETag>
    <Size>1024</Size>
    <StorageClass>STANDARD</StorageClass>
  </Contents>
  <Contents>
    <Key>notes.txt</Key>
    <LastModified>2024-01-02T00:00:00.000Z</LastModified>
    <ETag>"def456"</ETag>
    <Size>42</Size>
    <StorageClass>STANDARD</StorageClass>
  </Contents>
</ListBucketResult>`

	results, err := parseListBucketResult([]byte(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if results.Name != "example-bucket" {
		t.Errorf("Name = %q, want %q", results.Name, "example-bucket")
	}

	if len(results.Contents) != 2 {
		t.Fatalf("len(Contents) = %d, want 2", len(results.Contents))
	}

	if got, want := results.Contents[0].Key, "photo.jpg"; got != want {
		t.Errorf("Contents[0].Key = %q, want %q", got, want)
	}
	if got, want := results.Contents[0].Size, "1024"; got != want {
		t.Errorf("Contents[0].Size = %q, want %q", got, want)
	}
	if got, want := results.Contents[1].Key, "notes.txt"; got != want {
		t.Errorf("Contents[1].Key = %q, want %q", got, want)
	}
}

func TestParseListBucketResultEmpty(t *testing.T) {
	const body = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Name>empty-bucket</Name>
  <IsTruncated>false</IsTruncated>
</ListBucketResult>`

	results, err := parseListBucketResult([]byte(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results.Contents) != 0 {
		t.Errorf("len(Contents) = %d, want 0", len(results.Contents))
	}
}

func TestParseListBucketResultInvalid(t *testing.T) {
	if _, err := parseListBucketResult([]byte("not xml")); err == nil {
		t.Error("expected error for invalid XML, got nil")
	}
}
