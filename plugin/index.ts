const matchRe = RegExp('https://etf2l.org/matches/(\\d+)/')
const api_url = "http://localhost:8080/etf2l_match_page/"

class Log {
    id: number
    match_map: string
    uploaded_at: Date
    constructor(data: Object) {
        this.id = data["id"]
        this.match_map = data["match_map"]
        this.uploaded_at = new Date(data["uploaded_at"] * 1000)
    }
}

function getMatchID(): number {
    let url = document.URL;

    let match = url.match(matchRe);

    if (match === null || match.length < 1) {
        return;
    }
    return parseInt(match[0])
}

function getLogsFromAPI(match_id: number): Array<Log> {
    const apiResponse = {
        "logs": [
            {
                "id": 3137135,
                "match_map": "cp_sunshine",
                "uploaded_at": 	1645608660,
            },
            {
                "id": 3137108,
                "match_map": "cp_granary_pro_rc8",
                "uploaded_at": 	1645606860,
            },
        ]
    }

    let logs: Array<Log> = []

    for (let logData of apiResponse.logs) {
        const l = new Log(logData)
        logs.push(l)
    }
    return logs
}

function addLogLinks() {
    const match_id = getMatchID()
    const logs = getLogsFromAPI(match_id)

    let LogList = document.createElement("ul");

    for (const log of logs) {
        let logItem = document.createElement("li");
        logItem.innerHTML =`<a href="https://logs.tf/${log.id}"> ${log.match_map} | ${log.uploaded_at.toLocaleString()} </a>`;

        LogList.appendChild(logItem);
    }
    let LogHeader = document.createElement("div");
    LogHeader.className = "offi"
    LogHeader.innerHTML = `<h2>${logs.length} Logs</h2>`;

    LogHeader.append(LogList)

    let playersSection = document.getElementsByClassName("fix match-players");
    if (playersSection === null || playersSection.length < 1) {
        return;
    }
    playersSection[0].after(LogHeader)
}

addLogLinks()