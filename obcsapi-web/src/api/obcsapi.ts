
let host:string = localStorage.getItem("host")||"http://localhost:8900";

export const fetchData = async () => {
  const response = await fetch('https://jsonplaceholder.typicode.com/todos/1');
  return response.json();
};


export const ResetServerHost = () => {
  host = localStorage.getItem("host")||"http://localhost:8900";
};

export const GetObcsapiServerInfo = async () => {
  host = localStorage.getItem("host")||"http://localhost:8900";
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

export const ObcsapiGetMemos = async () => {
  const response = await fetch(host + '/api/v1/daily?day', {
    headers: {
      'Authorization': localStorage.getItem('token'),
    },
  });
  return response.json();
}

export const ObcsapiPostMemos = async (filekey: string, line: number, oldText: string, newText: string) => {
  const response = await fetch(host + '/api/v1/line', {
    method: 'POST',
    headers: {
      'Authorization': localStorage.getItem('token'),
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

export const ObcsapiTestJwt = async () => {
  const response = await fetch(host + '/api/v1/sayHello', {
    headers: {
      'Authorization': localStorage.getItem('token'),
    },
  });
  return response.text();
}
