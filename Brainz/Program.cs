// <copyright file="Program.cs" company="JP Dillingham">
//
//            ,. brainz
//      (¬º-°)¬ 
//
//      MIT License
//  
//      Copyright(c) 2019 JP Dillingham
//
//      Permission is hereby granted, free of charge, to any person obtaining a copy
//      of this software and associated documentation files (the "Software"), to deal
//      in the Software without restriction, including without limitation the rights
//      to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//      copies of the Software, and to permit persons to whom the Software is
//      furnished to do so, subject to the following conditions:
//
//      The above copyright notice and this permission notice shall be included in all
//      copies or substantial portions of the Software.
//
//      THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//      IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//      FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//      AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//      LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//      OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//      SOFTWARE.
// </copyright>

namespace brainz
{
    using Newtonsoft.Json;
    using System;
    using System.Linq;
    using System.Net.Http;
    using System.Threading.Tasks;
    using Utility.CommandLine;
    using System.Collections.Generic;
    using brainz.Responses;
    using brainz.Model;
    using System.Reflection;

    public enum Verbosity
    {
        None,
        Default,
        Verbose,
    }

    class Program
    {
        private static readonly string ApiRoot = @"https://musicbrainz.org/ws/2";

        private static readonly Assembly Assembly = Assembly.GetExecutingAssembly();
        private static readonly string AssemblyVersion = Assembly.GetCustomAttribute<AssemblyInformationalVersionAttribute>().InformationalVersion;
        private static readonly string AssemblyCompany = Assembly.GetCustomAttribute<AssemblyCompanyAttribute>().Company;
        private static readonly string AssemblyProduct = Assembly.GetCustomAttribute<AssemblyProductAttribute>().Product;
        private static readonly string UserAgent = $"{AssemblyProduct}/{AssemblyVersion} ({AssemblyCompany})";

        [Argument('l', "album", "The name of the album to search for.")]
        private static string Album { get; set; }

        [Argument('a', "artist", "The name of the artist to search for.")]
        private static string Artist { get; set; }

        [Argument('v', "verbosity", "The program output verbosity (None, Default, Verbose)")]
        private static string VerbosityLevel { get; set; }

        private static HttpClient Http { get; } = new HttpClient();
        private static Verbosity Verbosity { get; set; } = Verbosity.Default;

        private static Action<string> Output { get; } = (msg) => { if (Verbosity >= Verbosity.Default) Console.WriteLine(msg); };
        private static Action<string> Verbose { get; } = (msg) => { if (Verbosity == Verbosity.Verbose) Console.WriteLine(msg); };

        private static Func<string, string> ArtistRequest => (artist) => $"{ApiRoot}/artist/?query={artist}&fmt=json";
        private static Func<Guid, string> RecordingRequest => (mbid) => $"{ApiRoot}/release/{mbid}?inc=recordings&fmt=json";
        private static Func<Guid, int, int, string> ReleaseGroupRequest => (mbid, offset, limit) => $"{ApiRoot}/release-group?artist={mbid}&offset={offset}&limit={limit}&fmt=json";
        private static Func<Guid, int, int, string> ReleaseRequest => (mbid, offset, limit) => $"{ApiRoot}/release?release-group={mbid}&offset={offset}&limit={limit}&inc=media&fmt=json";

        private static async Task<T> Get<T>(string request)
        {
            var json = await Http.GetStringAsync(request).ConfigureAwait(false);
            return JsonConvert.DeserializeObject<T>(json);
        }

        private static async Task<Artist> GetBestArtist(string search)
        {
            Output($"Fetching artist matches for '{Artist}'...");

            var request = ArtistRequest(Artist);
            Verbose($"Fetching: {request}...");

            var artists = (await Get<ArtistResponse>(request)).Artists
                .OrderByDescending(a => a.Score);

            Verbose($"Fetched {artists.Count()} artists.");

            var bestArtist = artists.FirstOrDefault();

            foreach (var artist in Verbosity == Verbosity.Verbose ? artists : artists.Take(5))
            {
                var disambiguation = string.IsNullOrEmpty(artist.Disambiguation) ? string.Empty : $"({artist.Disambiguation})";
                Output($"{(artist.Id == bestArtist.Id ? "-->" : "   ")} {artist.Score.ToString().PadLeft(3)}%   {artist.Name} {disambiguation}");
            }

            Output($"Best artist match: {bestArtist.Name} (score: {bestArtist.Score}%)");

            return bestArtist;
        }

        static int Main(string[] args)
        {
            Output(string.Empty);
            Output("       ,. brainz");
            Output(" (¬º-°)¬ ");
            Output(string.Empty);

            Arguments.Populate();

            Artist = Artist ?? "atmosphere";
            Album = Album ?? "fishing";

            Enum.TryParse<Verbosity>(VerbosityLevel, out var verbosity);
            Verbosity = verbosity;

            Http.DefaultRequestHeaders.UserAgent.ParseAdd(UserAgent);

            return Task.Run(() => MainAsync(args)).ConfigureAwait(false).GetAwaiter().GetResult();
        }

        private static async Task<int> MainAsync(string[] args)
        {
            var bestArtist = await GetBestArtist(Artist);

            Console.WriteLine($"Fetching release group matches for artist '{bestArtist.Name}', album '{Album}'...");

            List<ReleaseGroup> releaseGroups = new List<ReleaseGroup>();
            ReleaseGroupResponse releaseGroupResponse = null;

            do
            {
                releaseGroupResponse = await Get<ReleaseGroupResponse>(ReleaseGroupRequest(bestArtist.Id, releaseGroups.Count, 100));
                releaseGroups.AddRange(releaseGroupResponse.ReleaseGroups);
            } while (releaseGroups.Count < releaseGroupResponse.ReleaseGroupCount);

            releaseGroups = releaseGroups
                .OrderByDescending(r => r.Title.SimilarityCaseInsensitive(Album))
                .ToList();

            var bestReleaseGroup = releaseGroups.FirstOrDefault();

            foreach (var releaseGroup in releaseGroups.Take(10))
            {
                Console.WriteLine($"{(releaseGroup.Title == bestReleaseGroup.Title ? "-->" : "   ")} {(releaseGroup.Title.SimilarityCaseInsensitive(Album) * 100).ToString("F0").PadLeft(3)}%   {releaseGroup.Title}");
            }

            Console.WriteLine();
            Console.WriteLine($"Best release group match: {bestReleaseGroup.Title} (score: {(bestReleaseGroup.Title.SimilarityCaseInsensitive(Album) * 100).ToString("F0")}%)");


            //Console.WriteLine(req);
            Console.WriteLine();
            Console.WriteLine($"Fetching releases for release group '{bestReleaseGroup.Title}'...");
            Console.WriteLine();

            List<Release> releases = new List<Release>();
            ReleaseResponse releasesResponse = null;

            string req;

            do
            {
                req = ReleaseRequest(bestReleaseGroup.Id, 0, 100);
                var releasesJson = await Http.GetStringAsync(req).ConfigureAwait(false);
                releasesResponse = JsonConvert.DeserializeObject<ReleaseResponse>(releasesJson);

                releases.AddRange(releasesResponse.Releases);
            } while (releases.Count < releasesResponse.ReleaseCount);
            
            releases = releases.OrderBy(r => r.Date.ToFuzzyDateTime()).ToList();

            //var bestRelease = releases
            //    .Where(r => r.Status == "Official")
            //    .Where(r => r.Country == "US")
            //    .Where(r => string.IsNullOrEmpty(r.Disambiguation)).FirstOrDefault();
            var bestRelease = releases.FirstOrDefault();

            foreach (var release in releases)
            {
                string trackStr = string.Join("+", release.Media.Select(m => m.TrackCount));
                string mediaStr = string.Join("+", release.Media.Select(m => m.Format));

                Console.WriteLine($"{(release.Id == bestRelease.Id ? "-->" : "   ")} {release.Title}\t{release.Country}\t{release.Date.PadRight(10)}\t{trackStr}\t{mediaStr}\t{release.Disambiguation}");
            }

            Console.WriteLine();
            Console.WriteLine($"Best release match: {bestRelease.Title}");

            req = RecordingRequest(bestRelease.Id);
            Console.WriteLine();
            Console.WriteLine($"Fetching recordings for release '{bestRelease.Title}' ({bestRelease.Id})");
            Console.WriteLine();

            var recordingJson = await Http.GetStringAsync(req).ConfigureAwait(false);
            var recordingResponse = JsonConvert.DeserializeObject<RecordingResponse>(recordingJson);

            foreach (var media in recordingResponse.Media)
            {
                foreach (var recording in media.Tracks.OrderBy(r => r.Position))
                {
                    Console.WriteLine($"{media.Position}{recording.Position.ToString("D2")} - {recording.Title}");
                }
            }

            return 0;
        }
    }
}
