/*
Copyright [2022] [IBM]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// cos_logger or Cloud object storage based logger, plays well with any logging
// package that can write to an io.Writer, including the standard library's log package.
// or as an extension to uber-go/zap logger.

package writer

import (
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
)

type Logger struct {
	Enabled bool
	buffer  *bytes.Buffer
	// When buffer reaches the size of MaxSize, it tries to sync with object store.
	MaxSize int64
	// Whether to compress before syncing the buffer.
	Compress bool
	// Current size of the buffer.
	size   int64
	mu     sync.Mutex
	Writer *Writer
}

// ensure we always implement io.WriteCloser
var _ io.WriteCloser = (*Logger)(nil)

func (l *Logger) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	writeLen := int64(len(p))
	if l.size+writeLen >= l.MaxSize {
		if err := l.syncBuffer(); err != nil {
			return 0, err
		}
	}
	if n, err = l.buffer.Write(p); err != nil {
		return n, err
	}
	l.size = l.size + int64(n)
	return n, nil
}

func (l *Logger) syncBuffer() error {
	var err error
	err = l.writeToObjectStore(l.Writer.DefaultBucketName,
		time.Now().Format(time.RFC3339Nano), l.buffer.Bytes())
	if err != nil {
		return err
	}
	l.buffer.Reset()
	l.size = 0
	return nil
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.syncBuffer()
}

func (l *Logger) load(o ObjectStoreConfig) error {
	cosCredentials := credentials.NewStaticCredentials(o.AccessKey, o.SecretKey, o.Token)
	// Create client config
	var conf = aws.NewConfig().
		WithRegion(o.Region).
		WithEndpoint(o.ServiceEndpoint).
		WithCredentials(cosCredentials).
		WithS3ForcePathStyle(o.S3ForcePathStyle)

	var sess = session.Must(session.NewSession())
	l.Writer = &Writer{
		DefaultBucketName: o.DefaultBucketName,
		client:            s3.New(sess, conf),
		mu:                sync.Mutex{},
	}
	input := &s3.CreateBucketInput{
		Bucket: aws.String(o.DefaultBucketName),
	}
	if o.CreateBucket {
		_, err := l.Writer.client.CreateBucket(input)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Logger) writeToObjectStore(bucketName string, key string, content []byte) error {
	input := s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(content),
	}

	_, err := l.Writer.client.PutObject(&input)
	return err
}

func (l *Logger) LoadDefaults(config ObjectStoreConfig) error {
	err := l.load(config)
	if err != nil {
		return err
	}
	if l.buffer == nil {
		l.buffer = new(bytes.Buffer)
	}
	return nil
}
