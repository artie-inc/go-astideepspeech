package astideepspeech

/*
#cgo LDFLAGS: -ldeepspeech
#include "deepspeech_wrap.h"
*/
import "C"
import "unsafe"

// Model represents a DeepSpeech model
type Model struct {
	beamWidth int
	modelPath string
	w         *C.ModelWrapper
}

// New creates a new Model
//
// modelPath          The path to the frozen model graph.
// beamWidth          The beam width used by the decoder. A larger beam width generates better results at the cost of decoding time.
func New(modelPath string, maxBatchSize, batchTimeoutMicros, numBatchThreads int) *Model {
	return &Model{
		modelPath: modelPath,
		w:         C.New(C.CString(modelPath), C.int(maxBatchSize), C.int(batchTimeoutMicros), C.int(numBatchThreads)),
	}
}

// Close closes the model properly
func (m *Model) Close() error {
	C.Close(m.w)
	return nil
}

// EnableExternalScorer enables decoding using an external scorer.
//
// scorerPath		The path to the scorer (kenlm language model) file
func (m *Model) EnableExternalScorer(scorerPath string) {
	C.EnableExternalScorer(m.w, C.CString(scorerPath))
}

// DisableExternalScorer disables decoding using an external scorer.
func (m *Model) DisableExternalScorer() {
	C.DisableExternalScorer(m.w)
}

// GetModelBeamWidth Get beam width value used by the model. If SetModelBeamWidth
// was not called before, will return the default value loaded from the model file.
func (m *Model) GetModelBeamWidth() uint {
	return uint(C.GetModelBeamWidth(m.w))
}

// SetModelBeamWidth  Set beam width value used by the model.
func (m *Model) SetModelBeamWidth(beamWidth uint) int {
	return int(C.SetModelBeamWidth(m.w, C.uint(m.beamWidth)))
}

// GetModelSampleRate read the sample rate that was used to produce the model file.
func (m *Model) GetModelSampleRate() int {
	return int(C.GetModelSampleRate(m.w))
}

// SetScorerAlphaBeta Set hyperparameters alpha and beta of the external scorer.
func (m *Model) SetScorerAlphaBeta(alpha, beta float32) int {
	return int(C.SetScorerAlphaBeta(m.w, C.float(alpha), C.float(beta)))
}

// sliceHeader represents a slice header
type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

// SpeechToText uses the DeepSpeech model to perform Speech-To-Text.
// buffer     A 16-bit, mono raw audio signal at the appropriate sample rate.
// bufferSize The number of samples in the audio signal.
func (m *Model) SpeechToText(buffer []int16, bufferSize uint) string {
	str := C.STT(m.w, (*C.short)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buffer)).Data)), C.uint(bufferSize))
	defer C.FreeString(str)
	retval := C.GoString(str)
	return retval
}

type TokenMetadata C.struct_TokenMetadata

func (tm *TokenMetadata) Text() string {
	return C.GoString(C.TokenMetadata_GetText((*C.TokenMetadata)(unsafe.Pointer(tm))))
}

func (tm *TokenMetadata) Timestep() uint {
	return uint(C.TokenMetadata_GetTimestep((*C.TokenMetadata)(unsafe.Pointer(tm))))
}

func (tm *TokenMetadata) StartTime() float32 {
	return float32(C.TokenMetadata_GetStartTime((*C.TokenMetadata)(unsafe.Pointer(tm))))
}

type CandidateTranscript C.struct_CandidateTranscript

func (ct *CandidateTranscript) Tokens() []TokenMetadata {

	// numItems := int32(C.Metadata_GetNumItems((*C.Metadata)(unsafe.Pointer(m))))
	// allItems := C.Metadata_GetItems((*C.Metadata)(unsafe.Pointer(m)))
	// return (*[1 << 30]MetadataItem)(unsafe.Pointer(allItems))[:numItems:numItems]
	numItems := int32(C.CandidateTranscript_GetNumTokens((*C.CandidateTranscript)(unsafe.Pointer(ct))))
	allItems := C.CandidateTranscript_GetTokens((*C.CandidateTranscript)(unsafe.Pointer(ct)))
	return (*[1 << 30]TokenMetadata)(unsafe.Pointer(allItems))[:numItems:numItems]
	// return C.CandidateTranscript_getTokens(ct)
}

func (ct *CandidateTranscript) NumTokens() uint {
	return uint(C.CandidateTranscript_GetNumTokens((*C.CandidateTranscript)(unsafe.Pointer(ct))))
}

func (ct *CandidateTranscript) Confidence() float64 {
	return float64(C.CandidateTranscript_GetConfidence((*C.CandidateTranscript)(unsafe.Pointer(ct))))
}

type Metadata C.struct_Metadata

func (m *Metadata) NumTranscripts() int32 {
	return int32(C.Metadata_GetNumTranscripts((*C.Metadata)(unsafe.Pointer(m))))
}

func (m *Metadata) Transcripts() []CandidateTranscript {
	numItems := int32(C.Metadata_GetNumTranscripts((*C.Metadata)(unsafe.Pointer(m))))
	allItems := C.Metadata_GetCandidateTranscripts((*C.Metadata)(unsafe.Pointer(m)))
	return (*[1 << 30]CandidateTranscript)(unsafe.Pointer(allItems))[:numItems:numItems]
}

// Close frees the Metadata structure properly
func (m *Metadata) Close() error {
	C.FreeMetadata((*C.Metadata)(unsafe.Pointer(m)))
	return nil
}

// --------------------------------------

// type MetadataItem C.struct_MetadataItem

// func (mi *MetadataItem) Character() string {
// 	return C.GoString(C.MetadataItem_GetCharacter((*C.MetadataItem)(unsafe.Pointer(mi))))
// }

// func (mi *MetadataItem) Timestep() int {
// 	return int(C.MetadataItem_GetTimestep((*C.MetadataItem)(unsafe.Pointer(mi))))
// }

// func (mi *MetadataItem) StartTime() float32 {
// 	return float32(C.MetadataItem_GetStartTime((*C.MetadataItem)(unsafe.Pointer(mi))))
// }

// // Metadata represents a DeepSpeech metadata output
// type Metadata C.struct_Metadata

// func (m *Metadata) NumItems() int32 {
// 	return int32(C.Metadata_GetNumItems((*C.Metadata)(unsafe.Pointer(m))))
// }

// func (m *Metadata) Confidence() float64 {
// 	return float64(C.Metadata_GetConfidence((*C.Metadata)(unsafe.Pointer(m))))
// }

// func (m *Metadata) Items() []MetadataItem {
// 	numItems := int32(C.Metadata_GetNumItems((*C.Metadata)(unsafe.Pointer(m))))
// 	allItems := C.Metadata_GetItems((*C.Metadata)(unsafe.Pointer(m)))
// 	return (*[1 << 30]MetadataItem)(unsafe.Pointer(allItems))[:numItems:numItems]
// }

// // Close frees the Metadata structure properly
// func (m *Metadata) Close() error {
// 	C.FreeMetadata((*C.Metadata)(unsafe.Pointer(m)))
// 	return nil
// }

//-------------------------------------------------------

// SpeechToTextWithMetadata uses the DeepSpeech model to perform Speech-To-Text.
// buffer     A 16-bit, mono raw audio signal at the appropriate sample rate.
// bufferSize The number of samples in the audio signal.
func (m *Model) SpeechToTextWithMetadata(buffer []int16, bufferSize uint, numResults uint) *Metadata {
	return (*Metadata)(unsafe.Pointer(C.STTWithMetadata(m.w, (*C.short)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buffer)).Data)), C.uint(bufferSize), C.uint(numResults))))
}

// Stream represent a streaming state
type Stream struct {
	sw *C.StreamWrapper
}

// CreateStream creates a new audio stream
//
// mw               The DeepSpeech model to use
func CreateStream(mw *Model) *Stream {
	return &Stream{
		sw: C.CreateStream(mw.w),
	}
}

// FeedAudioContent Feed audio samples to an ongoing streaming inference.
// aBuffer      An array of 16-bit, mono raw audio samples at the  appropriate sample rate.
// aBufferSize  The number of samples in @p aBuffer.
func (s *Stream) FeedAudioContent(buffer []int16, bufferSize uint) {
	C.FeedAudioContent(s.sw, (*C.short)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buffer)).Data)), C.uint(bufferSize))
}

func (s *Stream) FeedAudioContentFloat(buffer []float32, bufferSize uint) {
	C.FeedAudioContentFloat(s.sw, (*C.float)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buffer)).Data)), C.uint(bufferSize))
}

// IntermediateDecode Compute the intermediate decoding of an ongoing streaming inference.
// This is an expensive process as the decoder implementation isn't
// currently capable of streaming, so it always starts from the beginning
// of the audio.
func (s *Stream) IntermediateDecode() string {
	return C.GoString(C.IntermediateDecode(s.sw))
}

// IntermediateDecodeWithMetadata Compute the intermediate decoding of an ongoing streaming inference.
func (s *Stream) IntermediateDecodeWithMetadata(numResults uint) *Metadata {
	//return C.GoString(C.IntermediateDecodeWithMetadata(s.sw))
	return (*Metadata)(unsafe.Pointer(C.IntermediateDecodeWithMetadata(s.sw, C.uint(numResults))))
}

// FinishStream Signal the end of an audio signal to an ongoing streaming
// inference, returns the STT result over the whole audio signal.
func (s *Stream) FinishStream() string {
	str := C.FinishStream(s.sw)
	defer C.FreeString(str)
	retval := C.GoString(str)
	return retval
}

// FinishStreamWithMetadata Signal the end of an audio signal to an ongoing streaming
// inference, returns extended metadata.
func (s *Stream) FinishStreamWithMetadata(numResults uint) *Metadata {
	return (*Metadata)(unsafe.Pointer(C.FinishStreamWithMetadata(s.sw, C.uint(numResults))))
}

// DiscardStream Destroy a streaming state without decoding the computed logits.
// This can be used if you no longer need the result of an ongoing streaming
// inference and don't want to perform a costly decode operation.
func (s *Stream) FreeStream() {
	C.FreeStream(s.sw)
}

// PrintVersions Print version of this library and of the linked TensorFlow library.
// func PrintVersions() {
// 	C.PrintVersions()
// }
