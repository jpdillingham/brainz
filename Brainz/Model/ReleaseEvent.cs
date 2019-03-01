namespace Brainz.Model
{
    using Newtonsoft.Json;

    public class ReleaseEvent
    {
        [JsonProperty("date")]
        public string Date { get; set; }

        [JsonProperty("area")]
        public Area Area { get; set; }
    }
}
