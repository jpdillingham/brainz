namespace brainz.Model
{
    using Newtonsoft.Json;

    public class TextRepresentation
    {
        [JsonProperty("language")]
        public string Language { get; set; }

        [JsonProperty("script")]
        public string Script { get; set; }
    }
}
