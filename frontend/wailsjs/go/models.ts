export namespace db {
	
	export class APICredentials {
	    Provider: string;
	    Username: string;
	    Password: string;
	    APIKey: string;
	    BaseURL: string;
	    IsActive: boolean;
	    SearchByHash: boolean;
	    SearchByName: boolean;
	
	    static createFrom(source: any = {}) {
	        return new APICredentials(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Provider = source["Provider"];
	        this.Username = source["Username"];
	        this.Password = source["Password"];
	        this.APIKey = source["APIKey"];
	        this.BaseURL = source["BaseURL"];
	        this.IsActive = source["IsActive"];
	        this.SearchByHash = source["SearchByHash"];
	        this.SearchByName = source["SearchByName"];
	    }
	}

}

export namespace main {
	
	export class ErrorResult {
	    filename: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new ErrorResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filename = source["filename"];
	        this.reason = source["reason"];
	    }
	}
	export class SessionInfo {
	    found: boolean;
	    status: string;
	    root_path: string;
	    total_files: number;
	    done_files: number;
	
	    static createFrom(source: any = {}) {
	        return new SessionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.found = source["found"];
	        this.status = source["status"];
	        this.root_path = source["root_path"];
	        this.total_files = source["total_files"];
	        this.done_files = source["done_files"];
	    }
	}

}

