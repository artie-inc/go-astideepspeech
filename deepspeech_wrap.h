#ifdef __cplusplus
extern "C" {
#endif
    // typedef struct MetadataItem {
    //     char* character;
    //     int timestep;
    //     float start_time;
    // } MetadataItem;

    // typedef struct Metadata {
    //     MetadataItem* items;
    //     int num_items;
    //     double confidence;
    // } Metadata;

    typedef struct TokenMetadata {
        const char* const text;
        const unsigned int timestep;
        const float start_time;
    } TokenMetadata;

    typedef struct CandidateTranscript {
        const TokenMetadata* const tokens;
        const unsigned int num_tokens;
        const double confidence;
    } CandidateTranscript;

    typedef struct Metadata {
        const CandidateTranscript* const transcripts;
        const unsigned int num_transcripts;
    } Metadata;



    typedef void* ModelWrapper;
    ModelWrapper* New(const char* aModelPath, int maxBatchSize, int batchTimeoutMicros, int numBatchThreads);
    void Close(ModelWrapper* w);
    int EnableExternalScorer(ModelWrapper* w, const char* aScorerPath);
    int DisableExternalScorer(ModelWrapper* w);

    int SetScorerAlphaBeta(ModelWrapper* w, float aAlpha, float aBeta);
    unsigned int GetModelBeamWidth(ModelWrapper* w);
    unsigned int SetModelBeamWidth(ModelWrapper* w, unsigned int aBeamWidth);


    int GetModelSampleRate(ModelWrapper* w);

    char* STT(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize);
    Metadata* STTWithMetadata(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize, unsigned int aNumResults);

    typedef void* StreamWrapper;
    StreamWrapper* CreateStream(ModelWrapper* w);
    void FreeStream(StreamWrapper* sw);
    void FeedAudioContent(StreamWrapper* sw, const short* aBuffer, unsigned int aBufferSize);
    void FeedAudioContentFloat(StreamWrapper* sw, const float* aBuffer, unsigned int aBufferSize);
    char* IntermediateDecode(StreamWrapper* sw);
    Metadata* IntermediateDecodeWithMetadata(StreamWrapper* sw, unsigned int aNumResults);
    char* FinishStream(StreamWrapper* sw);
    Metadata* FinishStreamWithMetadata(StreamWrapper* sw, unsigned int aNumResults);

    int Metadata_GetNumTranscripts(Metadata* m);
    const CandidateTranscript* Metadata_GetCandidateTranscripts(Metadata* m) ;
    
    const TokenMetadata* CandidateTranscript_GetTokens(CandidateTranscript* ct);
    const unsigned int  CandidateTranscript_GetNumTokens(CandidateTranscript* ct);
    const double  CandidateTranscript_GetConfidence(CandidateTranscript* ct);

    const char* TokenMetadata_GetText(TokenMetadata* tm);
    const unsigned int TokenMetadata_GetTimestep(TokenMetadata* tm);
    const float TokenMetadata_GetStartTime(TokenMetadata* tm);


    void FreeString(char* s);
    void FreeMetadata(Metadata* m);
    // void PrintVersions();

    //char* DS_ErrorCodeToErrorMessage(int aErrorCode)

#ifdef __cplusplus
}
#endif
