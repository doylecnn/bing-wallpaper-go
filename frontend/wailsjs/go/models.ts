export namespace main {
	
	export class BingWallpaper {
	    url: string;
	    title: string;
	    copyright: string;
	    startdate: string;
	
	    static createFrom(source: any = {}) {
	        return new BingWallpaper(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.title = source["title"];
	        this.copyright = source["copyright"];
	        this.startdate = source["startdate"];
	    }
	}

}

