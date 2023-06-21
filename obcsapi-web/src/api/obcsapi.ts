
let host: string = localStorage.getItem("host") || "http://localhost:8900";

export const fetchData = async () => {
  const response = await fetch('https://jsonplaceholder.typicode.com/todos/1');
  return response.json();
};


export const ResetServerHost = () => {
  host = localStorage.getItem("host") || "http://localhost:8900";
};

export const GetObcsapiServerInfo = async () => {
  host = localStorage.getItem("host") || "http://localhost:8900";
  const response = await fetch(host + '/info');
  return response.json();
};

export const ObcsapiLogin = async (username: string, password: string) => {
  const response = await fetch(host + '/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ "username": username, "password": password })
  });
  return response.json();
}

export const ObcsapiGetMemos = async (day: number) => {
  const response = await fetch(host + '/api/v1/daily?day=' + day, {
    headers: {
      'Authorization': localStorage.getItem('token') || "",
    },
  });
  return response.json();
}

export const ObcsapiPostMemos = async (filekey: string, line: number, oldText: string, newText: string) => {
  const response = await fetch(host + '/api/v1/line', {
    method: 'POST',
    headers: {
      'Authorization': localStorage.getItem('token') || "",
    },
    body: JSON.stringify({
      "line_num": line,
      "day": filekey,
      "old": oldText,
      "content": newText
    })
  });
  return response.json();
}

export const ObcsapiConfigGet = async () => {
  const response = await fetch(host + '/api/v1/config', {
    headers: {
      'Authorization': localStorage.getItem('token') || "",
    }
  });
  return response.json();
}

export const ObcsapiMentionGet = async () => {
  const response = await fetch(host + '/api/v1/mention', {
    headers: {
      'Authorization': localStorage.getItem('token') || "",
    }
  });
  return response.json();
}

export const ObcsapiUpdateBdGet = async () => {
  const response = await fetch(host + '/api/v1/updatebd', {
    headers: {
      'Authorization': localStorage.getItem('token') || "",
    }
  });
  return response.json();
}

export const ObcsapiUpdateConfig = async () => {
  const response = await fetch(host + '/api/v1/updateconfig', {
    headers: {
      'Authorization': localStorage.getItem('token') || "",
    }
  });
  return response.json();
}




export const ObcsapiConfigPost = async (bodyObject: any) => {
  const response = await fetch(host + '/api/v1/config', {
    method: 'POST',
    headers: {
      'Authorization': localStorage.getItem('token') || "",
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(bodyObject)
  });
  return response.text();
}

export const ObcsapiTestJwt = async () => {
  const response = await fetch(host + '/api/v1/sayHello', {
    headers: {
      'Authorization': localStorage.getItem('token') || "",
    },
  });
  return response;
}

export const ObcsapiServerInfo = async () => {
  const response = await fetch(host + '/info');
  return response.json();
}

export const ObcsapiTestMail = async () => {
  const response = await fetch(host + '/api/v1/mailtest', {
    headers: {
      'Authorization': localStorage.getItem('token') || "",
    },
  });
  return response.text();
}

export const ObcsapiTalk = async (text: string) => {
  const response = await fetch(host + '/api/v1/talk', {
    method: 'POST',
    headers: {
      'Authorization': localStorage.getItem('token') || "",
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ "content": text }),
  });
  return response.text();
}

export const ObcsapiListFile = async () => {
  const response = await fetch(host + '/api/v1/list', {
    method: 'GET',
    headers: {
      'Authorization': localStorage.getItem('token') || "",
    },
  });
  return response.json();
}

export const ObcsapiTextGet = async (filekey: string) => {
  const response = await fetch(host + '/api/v1/text?fileKey=' + encodeURI(filekey), {
    method: 'GET',
    headers: {
      'Authorization': localStorage.getItem('token') || "",
    },
  });
  return response.text();
}

export const ObcsapiTextPost = async (filekey: string, text: string) => {
  const response = await fetch(host + '/api/v1/text?fileKey=' + encodeURI(filekey), {
    method: 'POST',
    headers: {
      'Authorization': localStorage.getItem('token') || "",
    },
    body: text,
  });
  return response;
}


export const ObcsapiUpdateCache = async (filekey: string) => {
  const response = await fetch(host + '/api/v1/cacheupdate?key=' + encodeURI(filekey), {
    method: 'POST',
    headers: {
      'Authorization': localStorage.getItem('token') || "",
    }
  });
  return response;
}

export const ObcsapiSerchKvCache = async (key: string) => {
  const response = await fetch(host + '/api/v1/search', {
    method: 'POST',
    headers: {
      'Authorization': localStorage.getItem('token') || "",
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ "key": key }),
  });
  return response;
}
