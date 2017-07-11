package scraper_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"strings"

	"github.com/gkats/scraper"
)

func TestScrape(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, resultsHTML)
	}))
	defer ts.Close()

	ads, err := scraper.Scrape(ts.URL)
	if err != nil {
		t.Error(err)
	}

	testCases := []struct {
		want interface{}
		got  interface{}
	}{
		{3, len(ads)},
		// First ad
		{"Women Shoes | Reebok GR", ads[0].H1},
		{"Reebok.com", ads[0].H2},
		{"www.reebok.com/Reebok+women+shoes", ads[0].Path},
		{"Flash Sale has started. Earn -25% more in Reebok products!", ads[0].Desc},
		{"<div class=\"ellip\">Free returns - Official Store</div>", ads[0].GetRest()},
		{1, ads[0].Position},
		// Second ad
		{"Reebok Women's", ads[1].H1},
		{"Latest Arrivals Are Here - cosmossport.gr", ads[1].H2},
		{"www.cosmossport.gr/reebok/womens", ads[1].Path},
		{"The latest arrivals in Reebok women's are at Cosmos Sport!", ads[1].Desc},
		{2, ads[1].Position},
		// Third ad
		{"New Νike Basketball Shoes", ads[2].H1},
		{"Catch the 10day Offer", ads[2].H2},
		{"www.zakcret.gr/nike/basket", ads[2].Path},
		{"New Releases in Nike Basketball Shoes. View the Collection, Buy Online!", ads[2].Desc},
		{3, ads[2].Position},
	}

	for i, tc := range testCases {
		if tc.got != tc.want {
			t.Errorf("(%v) Expected %v, got %v", i, tc.want, tc.got)
		}
	}

	if strings.Count(
		ads[2].GetRest(),
		"<div class=\"ellip\">Pay &amp; On delivery - 100% Change guarantee</div>",
	) != 1 {
		t.Errorf("Expected ad (2) rest to contain rest...")
	}

	for i, ad := range ads {
		if ad.GetRaw() == "" {
			t.Errorf("Expected ad (%v) raw not to be blank", i)
		}
	}
}

func TestScrapeSplitsH1AndH2Correctly(t *testing.T) {
	// TODO check what happens with the 30char limits
}

func TestNewURL(t *testing.T) {
	testCases := []struct {
		k    string
		want string
	}{
		{"searchterm", "https://www.google.com/search?q=searchterm"},
		{"search+term", "https://www.google.com/search?q=search+term"},
		{"searching for sugarman", "https://www.google.com/search?q=searching+for+sugarman"},
	}

	for i, tc := range testCases {
		if got := scraper.NewURL(tc.k); got != tc.want {
			t.Errorf("(%v) Expected %v, got %v", i, tc.want, got)
		}
	}
}

const resultsHTML = `
<!doctype html>
<head>
</head>
<body>
<div class="content">
<div data-jibp="h" data-jiis="uc" id="taw" style=""><style>.spell{font-size:18px}.spell_orig{font-size:15px}#mss p{margin:0;padding-top:5px}#tads h3,#tadsb h3,#mbEnd h3{font-size:18px !important;}.ads-ad:not(:first-child){border-top:1px solid transparent}#center_col ._Ak{color:#545454}.ads-ad{line-height:18px}#center_col ._Ak a._kBb{color:#609}#center_col ._Ak a:active{color:#dd4b39}#center_col ._Ak ._Bu,#center_col ._Ak ._Bu a{color:#808080}._x2b{margin-left:-16px;margin-right:-16px}._tve{border-top:1px solid #ebebeb;margin-right:-16px}._IYf{border-top:1px solid #ebebeb;margin-left:-16px;margin-right:-16px}.GLOBAL__wtaic{}.ads-ad{padding:11px 16px 11px 16px}#center_col ._Ak{position:relative;margin-left:0px;margin-right:0px}#tads{padding-top:2px;padding-bottom:5px;margin-top:-6px;margin-bottom:19px}#tadsb{padding-top:0;margin-top:-10px;margin-bottom:-11px;}.ads-ad{margin-bottom:5px}#center_col #tadsb._Ak{border-bottom:0}#center_col ._hM{margin:12px -17px 0 0}#center_col ._hM{font-weight:normal;font-size:13px;float:right}._hM span+span{margin-left:3px}._hM .g-bbll,._xQj{margin:0px 0px 0px 0px;padding:0px 0px 0px 0px}.ads-content-only{}.ads-ad-as{}.ads-ad-as{padding-top:9px;padding-right:16px;padding-bottom:11px;padding-left:16px}._uWj{padding-top:11px;padding-right:16px;padding-bottom:10px;padding-left:16px}._vWj{padding-top:11px;padding-right:16px;padding-bottom:10px;padding-left:0}._Tkg{display:-webkit-box;overflow:hidden;text-overflow:ellipsis;-webkit-box-orient:vertical;-webkit-line-clamp:2}.ads-visurl{color:#006621;white-space:nowrap;font-size:14px;}#center_col .ads-visurl cite{color:#006621;vertical-align:bottom}._WGk{display:inline-block;max-width:558px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.ads-visurl ._mB{margin-right:7px;margin-left:0px}._mB{background-color:#fff;border-radius:3px;color:#006621;display:inline-block;font-size:11px;border:1px solid #006621;padding:1px 3px 0 2px;line-height:11px;vertical-align:baseline}.action-menu,.action-menu-button,.action-menu-item,.action-menu-panel,.action-menu-toggled-item,.selected{}._Fmb,._Fmb:hover,._Fmb.selected,._Fmb.selected:hover{background-color:white;background-image:none;border:0;border-radius:0;box-shadow:0 0 0 0;cursor:pointer;filter:none;height:12px;min-width:0;padding:0;transition:none;-webkit-user-select:none;width:13px}.action-menu .mn-dwn-arw{border-color:#006621 transparent;margin-top:-4px;margin-left:3px;left:0;}.action-menu:hover .mn-dwn-arw{border-color:#00591E transparent}.action-menu{display:inline;margin:0 3px;position:relative;-webkit-user-select:none}.action-menu-panel{left:0;padding:0;right:auto;top:12px;visibility:hidden}.action-menu-item,.action-menu-toggled-item{cursor:pointer;-webkit-user-select:none}.action-menu-item:hover{background-color:#eee}.action-menu-button,.action-menu-item a.fl,.action-menu-toggled-item div{color:#333;display:block;padding:7px 18px;text-decoration:none;outline:0}._Ak .action-menu{line-height:0}._Ak .action-menu .mn-dwn-arw{border-color:#006621 transparent}._Ak .action-menu:hover .mn-dwn-arw{border-color:#00591E transparent}.ads-ad .action-menu .g-bbll{display:inline-block;height:12px;width:13px}.g-bbl-container{background-color:#fff;border:1px solid rgba(0,0,0,0.2);box-shadow:0 4px 16px rgba(0,0,0,0.2);color:#666;position:absolute;z-index:9120}.g-bbl-container.g-bbl-full{border-left-width:0;border-right-width:0;width:100%}.g-bbl-container.g-bbl-dark{background-color:#2d2d2d;border:1px solid rgba(0,0,0,0.5);color:#adadad;z-index:9100}.g-bbl-triangle{border-left-color:transparent;border-right-color:transparent;border-width:0 9.5px 9.5px 9.5px;width:0px;border-style:solid;border-top-color:transparent;height:0px;position:absolute;z-index:9121}.g-bbl-triangle.g-bbl-dark{z-index:9101}.g-bbl-triangle-bg{border-bottom-color:#bababa}.g-bbl-triangle-bg.g-bbl-dark{border-bottom-color:#0e0e0e}.g-bbl-triangle-fg{border-bottom-color:#fff;margin-left:-9px;margin-top:1px}.g-bbl-dark .g-bbl-triangle-fg{border-bottom-color:#2d2d2d}.g-bblc{display:none}._lBb{padding:16px}._zFc{padding-top:12px}._lBb a{text-decoration:none}._lBb a:hover{text-decoration:underline}._NU{margin-top:-2px;position:relative;top:2px}._uEc{color:#ef6c00}._pxg{background-repeat:repeat-x;display:inline-block;overflow:hidden;position:relative}._pxg span{background-repeat:repeat-x;display:block}._sxg ._pxg,._pxg._Esh{background-image:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABoAAAAaCAQAAAADQ4RFAAAA6klEQVR4AZXQMWsCMRiH8SAnQacODgpyg8rh1EEQHXS5xaUdXA5KRUHo+/2/wdN3aBNK34TEZ0rCD86/S/140ZydI9WrVo3etUrk+dJ8Hdog2qYO9YjW16ARD0R7MCpHC+SnRTk6BHQoR0NAg43WvP1LYsbrWh0tN6SwG+3v53n6ItLj//6nFfcsuLOyhphwSZILk/R6nUm6/OQzE83yaGeiXR5dTXTNoSmSaJpGWyQ0aBLaplGc/EijHePkKdTwRLQP5uFurifRnjQ2ahHtzBhHbKw3orU2OvHJEme01JeTjfZ4XCLPPp6+AYsy7RMdMSvnAAAAAElFTkSuQmCC)}._sxg ._pxg span,._pxg._Esh span{background-image:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABoAAAAaCAQAAAADQ4RFAAAA9klEQVR4AZXUoY7iUABG4SYYEgwYwhOsx4MlqUaAQ/AGMxqHIUHDC4DnATaMx7MORVAEh5vAtzUN7M69TXuOurc5SfuLJjH8ykzCJmJ++qgefWVWjJq+M5vVojEYV4s2YFMlqrmBm1r5qC+nXz5ayFmUj/7IOYajkd//uffO/sfzUdZJXZTlIs1fr2WrDFutf79p6KqIq2FoiLadGDvt+HoTISbFk3eF6BZHMyFmxdFBiENR1PEU4qkTj6ZeHDNfTONRPvnDUj1z6ZFPHovq7uCkJ7/rOYG7ejhKwVrD+23DGqThaOVsIAk4cLYKR3NN8b/T/HX6C7jRb/QEnjPPAAAAAElFTkSuQmCC)}._ayg{background:url(/images/nav_logo242.png) no-repeat -100px -260px;display:inline-block;height:13px;overflow:hidden;position:relative;top:1px;width:65px}._ayg span{background:url(/images/nav_logo242.png) no-repeat -100px -275px;display:block;height:13px;width:65px}._icr{display:-webkit-box;overflow:hidden;text-overflow:ellipsis;-webkit-box-orient:vertical}._mDc{-webkit-line-clamp:2}._jcr{-webkit-line-clamp:3}._hcr{-webkit-line-clamp:4}._tig{min-height:36px}.ads-creative b{color:#6a6a6a}._aes:last-child{margin-bottom:13px}._yEo>li+li:before{content:' · '}._r2b{color:#545454;font-size:small;margin-left:8px;white-space:nowrap;vertical-align:bottom}._J2b{position:relative;top:0.16em;display:inline-block;width:0.615em;height:1em;background-image:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAaBAMAAABMRsE0AAAAG1BMVEX///8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB4Gco9AAAACXRSTlMAETNmdyKIRFWlqMe0AAAAhUlEQVR4AT3KgQaCMRiF4feXKcAv3UAQgJABkDQBhq6gzSlAIAAhXXfHPjpg+94HWB9OW7yUpVaBnbybDxqrLBwc9y69puKWdYarGsVX9/5/OB1h4/RSm6esuz+6fOTzUmNPpiKvzxhJJmBkAkaDYBTEKIhREKMgGAWB7wMCBYHVm1iqwA8F3SwZzS7fmgAAAABJRU5ErkJggg==);background-size:100% 100%;background-repeat:no-repeat;border:0;margin-left:0.19em;margin-right:0.354em}#rcnt ._fEc a:hover{text-decoration:none}._Jpj{padding:8px;padding-top:6px;padding-bottom:14px;color:#444444}._Jpj td{padding:12px;padding-top:12px;padding-bottom:4px;vertical-align:top}._t9j{white-space:nowrap}._Jpj div{padding-top:4px}._Jpj div._hrj{padding-top:0px}._Jpj table{border-spacing:0px}._grj{font-weight:bold}._s9j{font-weight:normal;max-width:120px;overflow:hidden}._u9j{font-weight:normal;color:#689F38}._frj{font-weight:normal;color:#F47B00}._G2b{text-decoration:none;color:#808080}._G2b .mn-dwn-arw{position:relative;display:inline-block;margin-left:3px;margin-bottom:2px}._MEo{border-top:1px solid rgba(0,0,0,.12);padding:0;width:100%;font-size:14px}._OEo{margin-bottom:-11px;padding:11px 0 0 0px}._yEo{overflow:hidden;text-overflow:ellipsis;white-space:nowrap}._wEo,._uEo{margin:0 -13px -2px 0;padding:4px 0 3px 28px;width:418px}._uEo{width:604px}._vEo{margin:0 -13px -2px 0;padding:4px 0 3px 28px;width:614px}._wEo,._uEo{}._NEo{border-bottom:1px solid rgba(0,0,0,.12);padding:11px 0 0 0}._zEo{margin:0 0 -18px 0;padding:13px 0 1px 0px}._wEo>li,._vEo>li,._uEo li{box-sizing:border-box;box-sizing:border-box;display:inline-block;padding:0 13px 2px 0;vertical-align:top;width:50%}._uEo li{width:33.33%}._wEo li,._zEo li,._uEo li{overflow:hidden;text-overflow:ellipsis;white-space:nowrap}._zEo>li{padding:0 0 18px 0;font-size:16px}._vEo{padding-top:10px;padding-bottom:4px;margin-bottom:-15px}._vEo>li{padding-bottom:15px}._wEo>li,._vEo>li,._uEo li{line-height:inherit}._yEo>li{display:inline;margin:0;padding:0;line-height:inherit}._LEo{padding-left:16px;line-height:48px;margin-right:32px;color:#1a0dab}._KEo{padding-top:12px;float:right;margin-right:8px;color:rgba(0,0,0,.54)}#center_col #tads._Ak{margin-bottom:1px}</style><div></div><div style="padding:0 16px"><div class="med"><div class="_cy" id="msg_box" style="display:none"><p class="card-section _fbd"><span><span class="spell" id="srfm"></span>&nbsp;<a class="spell" id="srfl"></a><br></span><span id="sif"><span class="spell_orig" id="sifm"></span>&nbsp;<a class="spell_orig" id="sifl"></a><br></span></p></div></div></div>
	<div id="tvcap">
		<div class="_Ak c" id="tads" aria-label="Ads" role="region" data-ved="0ahUKEwjyh4a3yMfTAhVHvRoKHbV-C-8QGAgd">
			<h2 class="hd">Ads</h2>
			<h2 class="_hM"></h2>
			<ol>
				<li class="ads-ad" data-hveid="30">
					<h3>
						<a style="display:none" href="/aclk?sa=l&amp;ai=DChcSEwiv1ou3yMfTAhVWgbIKHSJCC-cYABAAGgJscg&amp;sig=AOD64_0EnZju6gHOfQnln0-YYMVFSNJTlA&amp;q=&amp;ved=0ahUKEwjyh4a3yMfTAhVHvRoKHbV-C-8Q0QwIHw&amp;adurl=" id="s0p1c0"></a><a href="http://www.reebok.com/gr/women-shoes" id="vs0p1c0" onmousedown="return google.arwt(this)" ontouchstart="return google.arwt(this)" data-preconnect-urls="http://clickserve.dartsearch.net/" jsl="$t t-zxXzjt1d4B0;$x 0;" class="r-iuOs0kzKssFE">
						Women Shoes | Reebok GR - Reebok.com
						</a>
					</h3>
					<div class="ads-visurl">
						<span class="_mB">Ad</span>
						<cite class="_WGk">
							www.reebok.com/Reebok+women+shoes
						</cite>&lrm;
						<g-bubble class="action-menu ab_ctl r-igw7vfh0IWuM" jsl="$t t-R7dwiTmE0C4;$x 0;"><a href="javascript:void(0)" data-theme="0" data-width="230" class="g-bbll" aria-haspopup="true" role="button" jsaction="r.saTe4DDW138" data-rtid="igw7vfh0IWuM" jsl="$x 1;" data-ved="0ahUKEwjyh4a3yMfTAhVHvRoKHbV-C-8QJwgg"><span class="mn-dwn-arw"></span></a><div class="g-bblc" data-ved="0ahUKEwjyh4a3yMfTAhVHvRoKHbV-C-8QKAgh"><div class="_lBb"><div>This ad is based on your current search terms.</div><div class="_zFc r-irtAgH6mFShI" jsl="$t t-h6cwAtrxkFI;$x 0;"><span>Visit Google’s <a href="javascript:void();" jsaction="r.8Na6VOGeTa8" data-rtid="irtAgH6mFShI" jsl="$x 10;" data-ved="0ahUKEwjyh4a3yMfTAhVHvRoKHbV-C-8QKQgi">Why This Ad page</a> to learn more or opt out.</span></div></div></div></g-bubble></div><div class="_Ond _Bu"><span class="_uEc">4.6</span> <g-review-stars><span class="_ayg" aria-label="Rated 4.5 out of 5,"><span style="width:59px"></span></span></g-review-stars> <a href="/shopping/seller?q=reebok.com&amp;hl=en&amp;sa=X&amp;ved=0ahUKEwjyh4a3yMfTAhVHvRoKHbV-C-8QwQYIJA">rating</a> for reebok.com
					</div>
					<div class="ellip ads-creative">
						Flash Sale has started. Earn -25% more in <b>Reebok</b> products!
					</div>
					<div class="ellip">Free returns - Official Store</div>
				</li>

				<li class="ads-ad" data-hveid="38">
					<h3>
						<a style="display:none" href="/aclk?sa=l&amp;ai=DChcSEwiv1ou3yMfTAhVWgbIKHSJCC-cYABABGgJscg&amp;sig=AOD64_25h8p9O-nJCZrKKRfsNLIMFfK7hg&amp;q=&amp;ved=0ahUKEwjyh4a3yMfTAhVHvRoKHbV-C-8Q0QwIJw&amp;adurl=" id="s0p2c0"></a><a href="https://www.cosmossport.gr/cosmos/el/catalog/brands/-b-q-r-b/reebok-sport-1329.html" id="vs0p2c0" onmousedown="return google.arwt(this)" ontouchstart="return google.arwt(this)" data-preconnect-urls="http://www.cosmossport.gr/" jsl="$t t-zxXzjt1d4B0;$x 0;" class="r-irANNgrtpclU">
							Reebok Women's - Latest Arrivals Are Here - cosmossport.gr&lrm;
						</a>
					</h3>
					<div class="ads-visurl">
						<span class="_mB">Ad</span>
						<cite class="_WGk" style="max-width:453px">
							www.cosmossport.gr/reebok/womens
						</cite>&lrm;
						<g-bubble class="action-menu ab_ctl r-iIsN0B2gtYP8" jsl="$t t-R7dwiTmE0C4;$x 0;"><a href="javascript:void(0)" data-theme="0" data-width="230" class="g-bbll" aria-haspopup="true" role="button" jsaction="r.saTe4DDW138" data-rtid="iIsN0B2gtYP8" jsl="$x 1;" data-ved="0ahUKEwjyh4a3yMfTAhVHvRoKHbV-C-8QJwgo"><span class="mn-dwn-arw"></span></a><div class="g-bblc" data-ved="0ahUKEwjyh4a3yMfTAhVHvRoKHbV-C-8QKAgp"><div class="_lBb"><div>This ad is based on your current search terms.</div><div class="_zFc r-iKX_9o9NWFn4" jsl="$t t-h6cwAtrxkFI;$x 0;"><span>Visit Google’s <a href="javascript:void();" jsaction="r.8Na6VOGeTa8" data-rtid="iKX_9o9NWFn4" jsl="$x 10;" data-ved="0ahUKEwjyh4a3yMfTAhVHvRoKHbV-C-8QKQgq">Why This Ad page</a> to learn more or opt out.</span></div></div></div></g-bubble><span class="_r2b" data-hveid="43">281 180 8808</span>
					</div>
					<div class="ellip ads-creative">
						The latest arrivals in <b>Reebok</b> women's are at Cosmos Sport!
					</div>
					<div class="ellip">Pay with e-banking&nbsp;·&nbsp;Free returns&nbsp;·&nbsp;Free delivery</div><div class="ellip"><a href="/aclk?sa=l&amp;ai=DChcSEwiv1ou3yMfTAhVWgbIKHSJCC-cYABADGgJscg&amp;sig=AOD64_0ABWNaL2L73hfGTYiXPZRAWi-y-A&amp;ctype=107&amp;q=&amp;ved=0ahUKEwjyh4a3yMfTAhVHvRoKHbV-C-8QmxAILg&amp;adurl=" aria-hidden="true"><span class="_J2b"></span><span class="_vnd">Εθνικής Αντιστάσεως 57, Περιστέρι</span></a>&lrm;<span aria-hidden="true"> - </span>&lrm;<span class="_fEc"><g-bubble jsl="$t t-R7dwiTmE0C4;$x 0;" class="r-irttIzkXNap8"><a href="javascript:void(0)" data-theme="0" data-width="-2" class="g-bbll" aria-haspopup="true" role="button" jsaction="r.saTe4DDW138" data-rtid="irttIzkXNap8" jsl="$x 1;" data-ved="0ahUKEwjyh4a3yMfTAhVHvRoKHbV-C-8Q_kAILw"><span class="_G2b"><span>Open today · 9:00 AM – 9:00 PM</span><span class="mn-dwn-arw"></span></span></a><div class="g-bblc" data-ved="0ahUKEwjyh4a3yMfTAhVHvRoKHbV-C-8Q_UAIMA"><div class="_Jpj"><table role="presentation"><tbody><tr class="_grj"><td><div class="_hrj">Friday</div></td><td class="_t9j"><div class="_hrj">9:00 AM – 9:00 PM</div></td></tr><tr><td><div class="_hrj">Saturday</div></td><td class="_t9j"><div class="_hrj">9:00 AM – 8:00 PM</div></td></tr><tr><td><div class="_hrj">Sunday</div></td><td class="_t9j"><div class="_hrj">Closed</div></td></tr><tr><td><div class="_hrj">Monday</div><div class="_s9j">(International Workers' Day)</div></td><td class="_t9j"><div class="_hrj">9:00 AM – 9:00 PM</div><div class="_frj">Hours might differ</div></td></tr><tr><td><div class="_hrj">Tuesday</div></td><td class="_t9j"><div class="_hrj">9:00 AM – 9:00 PM</div></td></tr><tr><td><div class="_hrj">Wednesday</div></td><td class="_t9j"><div class="_hrj">9:00 AM – 9:00 PM</div></td></tr><tr><td><div class="_hrj">Thursday</div></td><td class="_t9j"><div class="_hrj">9:00 AM – 9:00 PM</div></td></tr></tbody></table></div></div></g-bubble></span></div>
				</li>
			</ol>
		</div>
	</div><script jsl="$t t-IbpwSJ5oNyI;$x 0;" class="r-iT3aJ5gsk3FA"></script><form action="/settings/ads/preferences" method="post" jsl="$t t-y3Vq91bkxp8;$x 0;" class="r-iTPdYXcBTGIc"><input value="ADvV1IeN3T5FpWw36Y7pwOQDl4RHFUZNCToxNDkzMzk2NzUxNDgz" name="token" type="hidden"><input value="ChJyZWVib2sgd29tZW4gc2hvZXMaFggBEAEiEAjT8Mv5zAUQ-ebV2_MCIAAaFQgCEAIiDwj0-7_AuQUQyL6bolsgABoVCAQQAyIPCMqB7fayBRDwjqrZCCAAGhUIARAEIg8Iiqfgic4FEKLO3IwCIAAg8oeGt8jH0wIqDnd3dy5nb29nbGUuY29tMo0DaHR0cHM6Ly93d3cuZ29vZ2xlLmNvbS9zZWFyY2g_c2NsaWVudD1wc3ktYWImc2l0ZT0mc291cmNlPWhwJnE9cmVlYm9rK3dvbWVuK3Nob2VzJm9xPXJlZWJvayt3b21lbitzaG9lcyZnc19sPWhwLjMuLjBpMzBrMWowaTEwaTMwazFsMy4yMjMwLjEwNjg3LjAuMTA5NjAuMjMuMjEuMC4yLjIuMC4xNTAuMjQ1Ni4zajE4LjIxLjAuLi4uMC4uLjFjLjEuNjQucHN5LWFiLi4wLjIzLjIzOTYuLi4wajBpMTMxazFqMGkxM2sxajBpMTNpMzBrMWowaTEzaTEwaTMwazFqMGkyMmkzMGsxajBpMjJpMTBpMzBrMS5ZRXhpM2ZLZlRFTSZwYng9MSZiYXY9b24uMixvci4mZnA9MSZiaXc9MTkyMCZiaWg9OTg3JmRwcj0xJnRjaD0xJmVjaD0xJnBzaT1BMjBEV2RhZE1hT0g2QVRzb2Fqd0RnLjE0OTMzOTY3NDA0NTkuMzohUFJFRj1JRD0wMDAwMDAwMDAwMDAwMDAwOkZGPTA6Vj0x" name="reasons" type="hidden"><input value="en" name="hl" type="hidden"></form>
</div>
<div class="med" id="res" role="main">
Search results
</div>
<div data-jibp="h" data-jiis="uc" id="bottomads" style=""><style>#tads h3,#tadsb h3,#mbEnd h3{font-size:18px !important;}.ads-ad:not(:first-child){border-top:1px solid transparent}#center_col ._Ak{color:#545454}.ads-ad{line-height:18px}#center_col ._Ak a._kBb{color:#609}#center_col ._Ak a:active{color:#dd4b39}#center_col ._Ak ._Bu,#center_col ._Ak ._Bu a{color:#808080}._x2b{margin-left:-16px;margin-right:-16px}._tve{border-top:1px solid #ebebeb;margin-right:-16px}._IYf{border-top:1px solid #ebebeb;margin-left:-16px;margin-right:-16px}.GLOBAL__wtaic{}.ads-ad{padding:11px 16px 11px 16px}#center_col ._Ak{position:relative;margin-left:0px;margin-right:0px}#tads{padding-top:2px;padding-bottom:5px;margin-top:-6px;margin-bottom:19px}#tadsb{padding-top:0;margin-top:-10px;margin-bottom:-11px;}.ads-ad{margin-bottom:5px}#center_col #tadsb._Ak{border-bottom:0}#center_col ._hM{margin:12px -17px 0 0}#center_col ._hM{font-weight:normal;font-size:13px;float:right}._hM span+span{margin-left:3px}._hM .g-bbll,._xQj{margin:0px 0px 0px 0px;padding:0px 0px 0px 0px}.ads-content-only{}.ads-ad-as{}.ads-ad-as{padding-top:9px;padding-right:16px;padding-bottom:11px;padding-left:16px}._uWj{padding-top:11px;padding-right:16px;padding-bottom:10px;padding-left:16px}._vWj{padding-top:11px;padding-right:16px;padding-bottom:10px;padding-left:0}._Tkg{display:-webkit-box;overflow:hidden;text-overflow:ellipsis;-webkit-box-orient:vertical;-webkit-line-clamp:2}.ads-visurl{color:#006621;white-space:nowrap;font-size:14px;}#center_col .ads-visurl cite{color:#006621;vertical-align:bottom}._WGk{display:inline-block;max-width:558px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.ads-visurl ._mB{margin-right:7px;margin-left:0px}._mB{background-color:#fff;border-radius:3px;color:#006621;display:inline-block;font-size:11px;border:1px solid #006621;padding:1px 3px 0 2px;line-height:11px;vertical-align:baseline}.action-menu,.action-menu-button,.action-menu-item,.action-menu-panel,.action-menu-toggled-item,.selected{}._Fmb,._Fmb:hover,._Fmb.selected,._Fmb.selected:hover{background-color:white;background-image:none;border:0;border-radius:0;box-shadow:0 0 0 0;cursor:pointer;filter:none;height:12px;min-width:0;padding:0;transition:none;-webkit-user-select:none;width:13px}.action-menu .mn-dwn-arw{border-color:#006621 transparent;margin-top:-4px;margin-left:3px;left:0;}.action-menu:hover .mn-dwn-arw{border-color:#00591E transparent}.action-menu{display:inline;margin:0 3px;position:relative;-webkit-user-select:none}.action-menu-panel{left:0;padding:0;right:auto;top:12px;visibility:hidden}.action-menu-item,.action-menu-toggled-item{cursor:pointer;-webkit-user-select:none}.action-menu-item:hover{background-color:#eee}.action-menu-button,.action-menu-item a.fl,.action-menu-toggled-item div{color:#333;display:block;padding:7px 18px;text-decoration:none;outline:0}._Ak .action-menu{line-height:0}._Ak .action-menu .mn-dwn-arw{border-color:#006621 transparent}._Ak .action-menu:hover .mn-dwn-arw{border-color:#00591E transparent}.ads-ad .action-menu .g-bbll{display:inline-block;height:12px;width:13px}.g-bbl-container{background-color:#fff;border:1px solid rgba(0,0,0,0.2);box-shadow:0 4px 16px rgba(0,0,0,0.2);color:#666;position:absolute;z-index:9120}.g-bbl-container.g-bbl-full{border-left-width:0;border-right-width:0;width:100%}.g-bbl-container.g-bbl-dark{background-color:#2d2d2d;border:1px solid rgba(0,0,0,0.5);color:#adadad;z-index:9100}.g-bbl-triangle{border-left-color:transparent;border-right-color:transparent;border-width:0 9.5px 9.5px 9.5px;width:0px;border-style:solid;border-top-color:transparent;height:0px;position:absolute;z-index:9121}.g-bbl-triangle.g-bbl-dark{z-index:9101}.g-bbl-triangle-bg{border-bottom-color:#bababa}.g-bbl-triangle-bg.g-bbl-dark{border-bottom-color:#0e0e0e}.g-bbl-triangle-fg{border-bottom-color:#fff;margin-left:-9px;margin-top:1px}.g-bbl-dark .g-bbl-triangle-fg{border-bottom-color:#2d2d2d}.g-bblc{display:none}._lBb{padding:16px}._zFc{padding-top:12px}._lBb a{text-decoration:none}._lBb a:hover{text-decoration:underline}._NU{margin-top:-2px;position:relative;top:2px}._icr{display:-webkit-box;overflow:hidden;text-overflow:ellipsis;-webkit-box-orient:vertical}._mDc{-webkit-line-clamp:2}._jcr{-webkit-line-clamp:3}._hcr{-webkit-line-clamp:4}._tig{min-height:36px}.ads-creative b{color:#6a6a6a}._aes:last-child{margin-bottom:13px}._Eqs{display:-webkit-box;-webkit-box-orient:vertical;-webkit-line-clamp:2;overflow:hidden;text-overflow:ellipsis}._yEo>li+li:before{content:' · '}._MEo{border-top:1px solid rgba(0,0,0,.12);padding:0;width:100%;font-size:14px}._OEo{margin-bottom:-11px;padding:11px 0 0 0px}._yEo{overflow:hidden;text-overflow:ellipsis;white-space:nowrap}._wEo,._uEo{margin:0 -13px -2px 0;padding:4px 0 3px 28px;width:418px}._uEo{width:604px}._vEo{margin:0 -13px -2px 0;padding:4px 0 3px 28px;width:614px}._wEo,._uEo{}._NEo{border-bottom:1px solid rgba(0,0,0,.12);padding:11px 0 0 0}._zEo{margin:0 0 -18px 0;padding:13px 0 1px 0px}._wEo>li,._vEo>li,._uEo li{box-sizing:border-box;box-sizing:border-box;display:inline-block;padding:0 13px 2px 0;vertical-align:top;width:50%}._uEo li{width:33.33%}._wEo li,._zEo li,._uEo li{overflow:hidden;text-overflow:ellipsis;white-space:nowrap}._zEo>li{padding:0 0 18px 0;font-size:16px}._vEo{padding-top:10px;padding-bottom:4px;margin-bottom:-15px}._vEo>li{padding-bottom:15px}._wEo>li,._vEo>li,._uEo li{line-height:inherit}._yEo>li{display:inline;margin:0;padding:0;line-height:inherit}._LEo{padding-left:16px;line-height:48px;margin-right:32px;color:#1a0dab}._KEo{padding-top:12px;float:right;margin-right:8px;color:rgba(0,0,0,.54)}._J2b{position:relative;top:0.16em;display:inline-block;width:0.615em;height:1em;background-image:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAaBAMAAABMRsE0AAAAG1BMVEX///8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB4Gco9AAAACXRSTlMAETNmdyKIRFWlqMe0AAAAhUlEQVR4AT3KgQaCMRiF4feXKcAv3UAQgJABkDQBhq6gzSlAIAAhXXfHPjpg+94HWB9OW7yUpVaBnbybDxqrLBwc9y69puKWdYarGsVX9/5/OB1h4/RSm6esuz+6fOTzUmNPpiKvzxhJJmBkAkaDYBTEKIhREKMgGAWB7wMCBYHVm1iqwA8F3SwZzS7fmgAAAABJRU5ErkJggg==);background-size:100% 100%;background-repeat:no-repeat;border:0;margin-left:0.19em;margin-right:0.354em}#rcnt ._fEc a:hover{text-decoration:none}._Jpj{padding:8px;padding-top:6px;padding-bottom:14px;color:#444444}._Jpj td{padding:12px;padding-top:12px;padding-bottom:4px;vertical-align:top}._t9j{white-space:nowrap}._Jpj div{padding-top:4px}._Jpj div._hrj{padding-top:0px}._Jpj table{border-spacing:0px}._grj{font-weight:bold}._s9j{font-weight:normal;max-width:120px;overflow:hidden}._u9j{font-weight:normal;color:#689F38}._frj{font-weight:normal;color:#F47B00}._G2b{text-decoration:none;color:#808080}._G2b .mn-dwn-arw{position:relative;display:inline-block;margin-left:3px;margin-bottom:2px}</style>
	<div class="_Ak c" id="tadsb" aria-label="Ads" role="region" data-ved="0ahUKEwjLyazHysfTAhXMfxoKHfQ4AKIQ9AkIdQ">
		<h2 class="hd">Ads</h2>
		<h2 class="_hM"></h2>
		<ol>
			<li class="ads-ad" data-hveid="118">
				<h3>
					<a style="display:none" href="/aclk?sa=l&amp;ai=DChcSEwjE8bDHysfTAhWqCtMKHVpHD_IYABAKGgJ3Yg&amp;sig=AOD64_3jzMN7NR-FHSQ-on_vmBlI2yI8SA&amp;q=&amp;ved=0ahUKEwjLyazHysfTAhXMfxoKHfQ4AKIQ0QwIdw&amp;adurl=" id="s3p1c0"></a><a href="/aclk?sa=l&amp;ai=DChcSEwjE8bDHysfTAhWqCtMKHVpHD_IYABAKGgJ3Yg&amp;sig=AOD64_3jzMN7NR-FHSQ-on_vmBlI2yI8SA&amp;adurl=&amp;q=" id="vs3p1c0" onmousedown="return google.arwt(this)" ontouchstart="return google.arwt(this)" data-preconnect-urls="http://www.zakcret.gr/" jsl="$t t-zxXzjt1d4B0;$x 0;" class="r-ifIWiUiodDdE">
						New Νike Basketball Shoes - Catch the 10day Offer&lrm;
					</a>
				</h3>
				<div class="ads-visurl">
					<span class="_mB">Ad</span><cite class="_WGk">www.zakcret.gr/nike/basket</cite>&lrm;<g-bubble class="action-menu ab_ctl r-ilGmjPO_oTjs" jsl="$t t-R7dwiTmE0C4;$x 0;"><a href="javascript:void(0)" data-theme="0" data-width="230" class="g-bbll" aria-haspopup="true" role="button" jsaction="r.saTe4DDW138" data-rtid="ilGmjPO_oTjs" jsl="$x 1;" data-ved="0ahUKEwjLyazHysfTAhXMfxoKHfQ4AKIQJwh4"><span class="mn-dwn-arw"></span></a><div class="g-bblc" data-ved="0ahUKEwjLyazHysfTAhXMfxoKHfQ4AKIQKAh5"><div class="_lBb"><div>This ad is based on your current search terms.</div><div class="_zFc r-iq1Ug_y_uV8o" jsl="$t t-h6cwAtrxkFI;$x 0;"><span>Visit Google’s <a href="javascript:void();" jsaction="r.8Na6VOGeTa8" data-rtid="iq1Ug_y_uV8o" jsl="$x 10;" data-ved="0ahUKEwjLyazHysfTAhXMfxoKHfQ4AKIQKQh6">Why This Ad page</a> to learn more or opt out.</span></div></div></div></g-bubble>
				</div>
				<div class="ellip ads-creative">
					New Releases in <b>Nike</b> Basketball Shoes. View the Collection, Buy Online!
				</div>
				<div class="ellip">Pay &amp; On delivery - 100% Change guarantee</div><ul class="_yEo">
			<li><a style="display:none" href="/aclk?sa=l&amp;ai=DChcSEwjE8bDHysfTAhWqCtMKHVpHD_IYABALGgJ3Yg&amp;sig=AOD64_1qLq936EKV9lH1_kPDWpoWvEeY3g&amp;q=&amp;ved=0ahUKEwjLyazHysfTAhXMfxoKHfQ4AKIQpigIfSgA&amp;adurl=" id="ads-3-0-1-1-0"></a><a href="http://www.zakcret.gr/eshop/%CE%95%CF%84%CE%B1%CE%B9%CF%81%CE%B5%CE%AF%CE%B5%CF%82" id="vads-3-0-1-1-0" onmousedown="return google.arwt(this)">Δείτε Όλα τα Brands</a></li>
		</ol>
	</div>
</div>
</div>
</body>
`
