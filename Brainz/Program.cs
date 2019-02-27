using Newtonsoft.Json;
using System;
using System.Linq;
using System.Net.Http;
using System.Threading.Tasks;
using Utility.CommandLine;
using Brainz.Artist;
using Brainz.ReleaseGroup;
using Brainz.Release;
using Brainz.Recording;

namespace Brainz
{
    class Program
    {
        [Argument('a', "artist", "The name of the artist to search for.")]
        private static string Artist { get; set; } = "wu-tang clan";

        [Argument('r', "release", "The name of the release (album) to search for.")]
        private static string Album { get; set; } = "wu-tang forever";

        private static HttpClient Http { get; } = new HttpClient();

        private static string ApiRoot = @"https://musicbrainz.org/ws/2";
        private static Func<string, string> ArtistRequest => (artist) => $"{ApiRoot}/artist/?query={artist}&fmt=json";
        private static Func<Guid, string> ReleaseGroupRequest => (mbid) => $"{ApiRoot}/artist/{mbid}?inc=release-groups&fmt=json";
        private static Func<Guid, string> ReleaseRequest => (mbid) => $"{ApiRoot}/release-group/{mbid}?inc=releases&fmt=json";
        private static Func<Guid, string> RecordingRequest => (mbid) => $"{ApiRoot}/release/{mbid}?inc=recordings&fmt=json";

        static void Main(string[] args)
        {
            Http.DefaultRequestHeaders.UserAgent.ParseAdd("Brainz/1.00 ( https://github.com/jpdillingham/Brainz )");
            Task.Run(() => MainAsync(args)).ConfigureAwait(false).GetAwaiter().GetResult();

            Console.ReadKey();
        }

        private static async void MainAsync(string[] args)
        {
            var req = ArtistRequest(Artist);
            Console.WriteLine($"Fetching {req}...");
            var artistJson = await Http.GetStringAsync(req).ConfigureAwait(false);
            var artistResponse = JsonConvert.DeserializeObject<ArtistResponse>(artistJson);
            Console.WriteLine($"{artistResponse.Artists.Count()} artists returned.");

            var bestArtist = artistResponse.Artists.OrderByDescending(a => a.Score).FirstOrDefault();

            Console.WriteLine($"Best artist match: {bestArtist.Name} (score: {bestArtist.Score})");

            req = ReleaseGroupRequest(bestArtist.Id);
            Console.WriteLine($"Fetching {req}...");
            var releaseGroupJson = await Http.GetStringAsync(ReleaseGroupRequest(bestArtist.Id)).ConfigureAwait(false);
            var releaseGroupResponse = JsonConvert.DeserializeObject<ReleaseGroupResponse>(releaseGroupJson);
            Console.WriteLine($"{releaseGroupResponse.ReleaseGroups.Count()} release groups returned.");

            var bestReleaseGroup = releaseGroupResponse.ReleaseGroups.FirstOrDefault();

            foreach (var releaseGroup in releaseGroupResponse.ReleaseGroups)
            {
                if (releaseGroup.Title.LevenshteinDistanceCaseInsensitive(Album) < bestReleaseGroup.Title.LevenshteinDistanceCaseInsensitive(Album))
                {
                    bestReleaseGroup = releaseGroup;
                }
            }

            Console.WriteLine($"Best release group match: {bestReleaseGroup.Title} (lev. distance: {bestReleaseGroup.Title.LevenshteinDistanceCaseInsensitive(Album)})");

            req = ReleaseRequest(bestReleaseGroup.Id);
            Console.WriteLine($"Fetching {req}...");
            var releasesJson = await Http.GetStringAsync(req).ConfigureAwait(false);
            var releaseResponse = JsonConvert.DeserializeObject<ReleaseResponse>(releasesJson);
            Console.WriteLine($"{releaseResponse.Releases} releases returned.");

            var bestRelease = releaseResponse.Releases
                .Where(r => r.Status == "Official")
                .Where(r => string.IsNullOrEmpty(r.Disambiguation)).FirstOrDefault();

            Console.WriteLine($"Best release match: {bestRelease.Title}");

            req = RecordingRequest(bestRelease.Id);
            Console.WriteLine($"Fetching {req}...");
            var recordingJson = await Http.GetStringAsync(req).ConfigureAwait(false);
            var recordingResponse = JsonConvert.DeserializeObject<RecordingResponse>(recordingJson);
            Console.WriteLine($"{recordingResponse.Media[0].Tracks.Count()} recordings returned.");

            foreach (var media in recordingResponse.Media)
            {
                foreach (var recording in media.Tracks.OrderBy(r => r.Position))
                {
                    Console.WriteLine($"{media.Position}{recording.Position.ToString("D2")} - {recording.Title}");
                }
            }
        }
    }
}
