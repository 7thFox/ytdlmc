package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

var (
	config    = flag.String("config", "~/.config/youtube-dl-multiconfig", "The path to your config file")
	simulate  = flag.Bool("simulate", false, "Print command and don't execute")
	tempFiles []*os.File
)

func main() {
	defer func() {
		for _, f := range tempFiles {
			if f != nil {
				os.Remove(f.Name())
			}
		}
	}()

	log.SetFlags(0)
	flag.Parse()

	file, err := ioutil.ReadFile(*config)
	if err != nil {
		log.Fatalf("Error opening file: %s\n", err.Error())
	}

	groups := make(map[string]ConfigGroup)
	if err := json.Unmarshal(file, &groups); err != nil {
		log.Fatalf("Error parsing file: %s\n", err.Error())
	}

	for name, group := range groups {
		log.Printf("Processing %s...\n", name)
		args := getArgs(group)
		if *simulate {
			log.Println(getCommandString(args))
		} else {
			cmd := exec.Command("youtube-dl", args...)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			time.Sleep(500) // "Can't read file" error otherwise

			err = cmd.Run()
			if err != nil {
				log.Fatalf("Error executing command: %s\nCommand: %s\n", err.Error(), getCommandString(args))
			}
		}
	}

	log.Println("Done.")

	fmt.Scanln()
}

func getCommandString(args []string) string {
	var sb strings.Builder
	sb.WriteString("youtube-dl")
	for _, arg := range args {
		sb.WriteRune(' ')
		sb.WriteString(arg)
	}
	return sb.String()
}

func getArgs(group ConfigGroup) []string {
	t := reflect.TypeOf(group)
	v := reflect.ValueOf(group)
	nFields := t.NumField()
	var args []string
	for i := 0; i < nFields; i++ {
		f := t.Field(i)
		if opt, ok := f.Tag.Lookup("option"); ok {
			val := v.Field(i)
			args = writeField(opt, val, args)
		}
	}

	return args
}

func writeField(opt string, fval reflect.Value, args []string) []string {
	switch fval.Kind() {
	case reflect.Bool:
		if fval.Bool() {
			return append(args, "--"+opt)
		}
	case reflect.Ptr:
		if !fval.IsNil() {
			pval := reflect.Indirect(fval)
			return writeField(opt, pval, args)
		}
	case reflect.String:
		return append(args, "--"+opt, fval.String())
	case reflect.Int:
		return append(args, "--"+opt, fmt.Sprintf("%d", fval.Int()))
	case reflect.Slice:
		log.Println(opt)
		if opt == "batch-file" {
			file, err := ioutil.TempFile(os.TempDir(), "youtube-dl-multiconfig-*")
			if err != nil {
				log.Fatalf("Failed to create temp file: %s\n", err.Error())
			}
			tempFiles = append(tempFiles, file)

			for i := 0; i < fval.Len(); i++ {
				if _, err = file.WriteString(fval.Index(i).String() + "\n"); err != nil {
					log.Fatalf("Failed to write to temp file: %s\n", err.Error())
				}
			}
			if err = file.Close(); err != nil {
				log.Fatalf("Failed to close temp file: %s\n", err.Error())
			}
			return append(args, "--"+opt, file.Name())
		}
		goto Err
	default:
		goto Err
	}
	return args

Err:
	log.Printf("Warning: Unhandled option type: %v\n", fval.Kind())
	return args
}

type ConfigGroup struct {
	// TODO JOSH
	BatchFile []string `json:"batch-file" option:"batch-file"`

	Help                       bool `json:"help" option:"help"`
	Version                    bool `json:"version" option:"version"`
	Update                     bool `json:"update" option:"update"`
	IgnoreErrors               bool `json:"ignore-errors" option:"ignore-errors"`
	AbortOnError               bool `json:"abort-on-error" option:"abort-on-error"`
	DumpUserAgent              bool `json:"dump-user-agent" option:"dump-user-agent"`
	ListExtractors             bool `json:"list-extractors" option:"list-extractors"`
	ExtractorDescriptions      bool `json:"extractor-descriptions" option:"extractor-descriptions"`
	ForceGenericExtractor      bool `json:"force-generic-extractor" option:"force-generic-extractor"`
	IgnoreConfig               bool `json:"ignore-config" option:"ignore-config"`
	FlatPlaylist               bool `json:"flat-playlist" option:"flat-playlist"`
	MarkWatched                bool `json:"mark-watched" option:"mark-watched"`
	NoMarkWatched              bool `json:"no-mark-watched" option:"no-mark-watched"`
	NoColor                    bool `json:"no-color" option:"no-color"`
	ForceIPV4                  bool `json:"force-ipv4" option:"force-ipv4"`
	ForceIPV6                  bool `json:"force-ipv6" option:"force-ipv6"`
	GeoBypass                  bool `json:"geo-bypass" option:"geo-bypass"`
	NoGeoBypass                bool `json:"no-geo-bypass" option:"no-geo-bypass"`
	PlaylistReverse            bool `json:"playlist-reverse" option:"playlist-reverse"`
	PlaylistRandom             bool `json:"playlist-random" option:"playlist-random"`
	XattrSetFilesize           bool `json:"xattr-set-filesize" option:"xattr-set-filesize"`
	HlsPreferNative            bool `json:"hls-prefer-native" option:"hls-prefer-native"`
	HlsPreferFfmpeg            bool `json:"hls-prefer-ffmpeg" option:"hls-prefer-ffmpeg"`
	HlsUseMpegts               bool `json:"hls-use-mpegts" option:"hls-use-mpegts"`
	RestrictFilenames          bool `json:"restrict-filenames" option:"restrict-filenames"`
	NoOverwrites               bool `json:"no-overwrites" option:"no-overwrites"`
	Continue                   bool `json:"continue" option:"continue"`
	NoContinue                 bool `json:"no-continue" option:"no-continue"`
	NoPart                     bool `json:"no-part" option:"no-part"`
	NoMtime                    bool `json:"no-mtime" option:"no-mtime"`
	WriteDescription           bool `json:"write-description" option:"write-description"`
	WriteInfoJSON              bool `json:"write-info-json" option:"write-info-json"`
	WriteAnnotations           bool `json:"write-annotations" option:"write-annotations"`
	NoCacheDir                 bool `json:"no-cache-dir" option:"no-cache-dir"`
	RmCacheDir                 bool `json:"rm-cache-dir" option:"rm-cache-dir"`
	WriteThumbnail             bool `json:"write-thumbnail" option:"write-thumbnail"`
	WriteAllThumbnails         bool `json:"write-all-thumbnails" option:"write-all-thumbnails"`
	ListThumbnails             bool `json:"list-thumbnails" option:"list-thumbnails"`
	Quiet                      bool `json:"quiet" option:"quiet"`
	NoWarnings                 bool `json:"no-warnings" option:"no-warnings"`
	Simulate                   bool `json:"simulate" option:"simulate"`
	SkipDownload               bool `json:"skip-download" option:"skip-download"`
	GetURL                     bool `json:"get-url" option:"get-url"`
	GetTitle                   bool `json:"get-title" option:"get-title"`
	GetID                      bool `json:"get-id" option:"get-id"`
	GetThumbnail               bool `json:"get-thumbnail" option:"get-thumbnail"`
	GetDescription             bool `json:"get-description" option:"get-description"`
	GetDuration                bool `json:"get-duration" option:"get-duration"`
	GetFilename                bool `json:"get-filename" option:"get-filename"`
	GetFormat                  bool `json:"get-format" option:"get-format"`
	DumpJSON                   bool `json:"dump-json" option:"dump-json"`
	DumpSingleJSON             bool `json:"dump-single-json" option:"dump-single-json"`
	PrintJSON                  bool `json:"print-json" option:"print-json"`
	Newline                    bool `json:"newline" option:"newline"`
	NoProgress                 bool `json:"no-progress" option:"no-progress"`
	ConsoleTitle               bool `json:"console-title" option:"console-title"`
	Verbose                    bool `json:"verbose" option:"verbose"`
	DumpPages                  bool `json:"dump-pages" option:"dump-pages"`
	WritePages                 bool `json:"write-pages" option:"write-pages"`
	PrintTraffic               bool `json:"print-traffic" option:"print-traffic"`
	CallHome                   bool `json:"call-home" option:"call-home"`
	NoCallHome                 bool `json:"no-call-home" option:"no-call-home"`
	NoPlaylist                 bool `json:"no-playlist" option:"no-playlist"`
	YesPlaylist                bool `json:"yes-playlist" option:"yes-playlist"`
	IncludeAds                 bool `json:"include-ads" option:"include-ads"`
	SkipUnavailableFragments   bool `json:"skip-unavailable-fragments" option:"skip-unavailable-fragments"`
	AbortOnUnavailableFragment bool `json:"abort-on-unavailable-fragment" option:"abort-on-unavailable-fragment"`
	KeepFragments              bool `json:"keep-fragments" option:"keep-fragments"`
	NoResizeBuffer             bool `json:"no-resize-buffer" option:"no-resize-buffer"`
	ID                         bool `json:"id" option:"id"`
	NoCheckCertificate         bool `json:"no-check-certificate" option:"no-check-certificate"`
	PreferInsecure             bool `json:"prefer-insecure" option:"prefer-insecure"`
	AddHeader                  bool `json:"add-header" option:"add-header"`
	BidiWorkaround             bool `json:"bidi-workaround" option:"bidi-workaround"`
	AllFormats                 bool `json:"all-formats" option:"all-formats"`
	PreferFreeFormats          bool `json:"prefer-free-formats" option:"prefer-free-formats"`
	ListFormats                bool `json:"list-formats" option:"list-formats"`
	YoutubeSkipDashManifest    bool `json:"youtube-skip-dash-manifest" option:"youtube-skip-dash-manifest"`
	WriteSub                   bool `json:"write-sub" option:"write-sub"`
	WriteAutoSub               bool `json:"write-auto-sub" option:"write-auto-sub"`
	AllSubs                    bool `json:"all-subs" option:"all-subs"`
	ListSubs                   bool `json:"list-subs" option:"list-subs"`
	Netrc                      bool `json:"netrc" option:"netrc"`
	ApListMso                  bool `json:"ap-list-mso" option:"ap-list-mso"`
	ExtractAudio               bool `json:"extract-audio" option:"extract-audio"`
	KeepVideo                  bool `json:"keep-video" option:"keep-video"`
	NoPostOverwrites           bool `json:"no-post-overwrites" option:"no-post-overwrites"`
	EmbedSubs                  bool `json:"embed-subs" option:"embed-subs"`
	EmbedThumbnail             bool `json:"embed-thumbnail" option:"embed-thumbnail"`
	AddMetadata                bool `json:"add-metadata" option:"add-metadata"`
	Xattrs                     bool `json:"xattrs" option:"xattrs"`
	PreferAvconv               bool `json:"prefer-avconv" option:"prefer-avconv"`
	PreferFfmpeg               bool `json:"prefer-ffmpeg" option:"prefer-ffmpeg"`

	DefaultSearch          *string `json:"default-search" option:"default-search"`
	ConfigLocation         *string `json:"config-location" option:"config-location"`
	Proxy                  *string `json:"proxy" option:"proxy"`
	SourceAddress          *string `json:"source-address" option:"source-address"`
	GeoVerificationProxy   *string `json:"geo-verification-proxy" option:"geo-verification-proxy"`
	GeoBypassCountry       *string `json:"geo-bypass-country" option:"geo-bypass-country"`
	GeoBypassIPBlock       *string `json:"geo-bypass-ip-block" option:"geo-bypass-ip-block"`
	PlaylistItems          *string `json:"playlist-items" option:"playlist-items"`
	MatchTitle             *string `json:"match-title" option:"match-title"`
	RejectTitle            *string `json:"reject-title" option:"reject-title"`
	MinFilesize            *string `json:"min-filesize" option:"min-filesize"`
	MaxFilesize            *string `json:"max-filesize" option:"max-filesize"`
	Date                   *string `json:"date" option:"date"`
	Datebefore             *string `json:"datebefore" option:"datebefore"`
	Dateafter              *string `json:"dateafter" option:"dateafter"`
	MatchFilter            *string `json:"match-filter" option:"match-filter"`
	DownloadArchive        *string `json:"download-archive" option:"download-archive"`
	LimitRate              *string `json:"limit-rate" option:"limit-rate"`
	BufferSize             *string `json:"buffer-size" option:"buffer-size"`
	HTTPChunkSize          *string `json:"http-chunk-size" option:"http-chunk-size"`
	ExternalDownloader     *string `json:"external-downloader" option:"external-downloader"`
	ExternalDownloaderArgs *string `json:"external-downloader-args" option:"external-downloader-args"`
	Output                 *string `json:"output" option:"output"`
	OutputNaPlaceholder    *string `json:"output-na-placeholder" option:"output-na-placeholder"`
	LoadInfoJSON           *string `json:"load-info-json" option:"load-info-json"`
	Cookies                *string `json:"cookies" option:"cookies"`
	CacheDir               *string `json:"cache-dir" option:"cache-dir"`
	Encoding               *string `json:"encoding" option:"encoding"`
	UserAgent              *string `json:"user-agent" option:"user-agent"`
	Referer                *string `json:"referer" option:"referer"`
	Format                 *string `json:"format" option:"format"`
	MergeOutputFormat      *string `json:"merge-output-format" option:"merge-output-format"`
	SubFormat              *string `json:"sub-format" option:"sub-format"`
	SubLang                *string `json:"sub-lang" option:"sub-lang"`
	Username               *string `json:"username" option:"username"`
	Password               *string `json:"password" option:"password"`
	Twofactor              *string `json:"twofactor" option:"twofactor"`
	VideoPassword          *string `json:"video-password" option:"video-password"`
	ApMso                  *string `json:"ap-mso" option:"ap-mso"`
	ApUsername             *string `json:"ap-username" option:"ap-username"`
	ApPassword             *string `json:"ap-password" option:"ap-password"`
	AudioFormat            *string `json:"audio-format" option:"audio-format"`
	AudioQuality           *string `json:"audio-quality" option:"audio-quality"`
	RecodeVideo            *string `json:"recode-video" option:"recode-video"`
	PostprocessorArgs      *string `json:"postprocessor-args" option:"postprocessor-args"`
	MetadataFromTitle      *string `json:"metadata-from-title" option:"metadata-from-title"`
	Fixup                  *string `json:"fixup" option:"fixup"`
	FfmpegLocation         *string `json:"ffmpeg-location" option:"ffmpeg-location"`
	Exec                   *string `json:"exec" option:"exec"`
	ConvertSubs            *string `json:"convert-subs" option:"convert-subs"`

	SocketTimeout    *int `json:"socket-timeout" option:"socket-timeout"`
	PlaylistStart    *int `json:"playlist-start" option:"playlist-start"`
	PlaylistEnd      *int `json:"playlist-end" option:"playlist-end"`
	MaxDownloads     *int `json:"max-downloads" option:"max-downloads"`
	MinViews         *int `json:"min-views" option:"min-views"`
	MaxViews         *int `json:"max-views" option:"max-views"`
	AgeLimit         *int `json:"age-limit" option:"age-limit"`
	Retries          *int `json:"retries" option:"retries"`
	FragmentRetries  *int `json:"fragment-retries" option:"fragment-retries"`
	AutonumberStart  *int `json:"autonumber-start" option:"autonumber-start"`
	SleepInterval    *int `json:"sleep-interval" option:"sleep-interval"`
	MaxSleepInterval *int `json:"max-sleep-interval" option:"max-sleep-interval"`
}
