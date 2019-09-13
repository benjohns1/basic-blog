import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

interface Token {
  token: string
}

interface Error {
  error: string
}

@Injectable({
  providedIn: 'root'
})
export class AuthenticationService {

  constructor(private router: Router, private http: HttpClient) { }

  public login(user: string, password: string): Observable<boolean|string> {
    const postData = JSON.stringify({ user, password });
    const success$ = new Observable<boolean|string>(obs => {
      this.http.post(`http://localhost:3000/api/v1/authenticate`, postData, { observe: 'response' }).subscribe(response => {
        if (response.status === 200) {
          const data = response.body as Token;
          if (data) {
            window.sessionStorage.setItem("token", (response.body as Token).token);
            obs.next(true);
            obs.complete();
            return
          }
        }
        obs.next("Unauthenticated");
        obs.complete();
      }, error => {
        obs.next((error.error as Error).error);
        obs.complete();
      });
    });
    return success$;
  }

  public logout() {
    window.sessionStorage.removeItem("token");
    this.router.navigate(['/']);
  }

  public isLoggedIn() {
    return !!window.sessionStorage.getItem("token");
  }
}
