import { Injectable } from '@angular/core';
import { Post } from './post';
import { NewComment } from './comment';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

export enum BlogPostFilter {
  None,
  Deleted,
}

@Injectable({
  providedIn: 'root'
})
export class BlogService {

  private baseUrl: string = "http://localhost:3000/api/v1";

  constructor(private http: HttpClient) { }

  public getPosts(postFilter: BlogPostFilter = BlogPostFilter.None): Observable<Post[]> {
    return this.http.get<Post[]>(`${this.baseUrl}/post/`, { headers: this.headers() }).pipe(
      map((response: Post[]) => {
        return response.filter(p => {
          switch (postFilter) {
            case BlogPostFilter.None:
              return !!p.deleted === false;
            case BlogPostFilter.Deleted:
              return p.deleted === true;
            default:
              return true;
          }
        });
      })
    );
  }

  public getPost(id: number): Observable<Post> {
    return this.http.get<Post>(`${this.baseUrl}/post/${id}`, { headers: this.headers() });
  }

  public newComment(comment: NewComment): Observable<Object> {
    if (!comment.postId || comment.postId <= 0) {
      return new Observable(o => {
        o.error("Post ID must be a positive integer");
      });
    }

    return this.http.post(`${this.baseUrl}/post/${comment.postId}/comment`, {
      commenter: comment.commenter,
      body: comment.body,
    }, { headers: this.headers() });
  }

  public newPost(post: Post): Observable<Post> {
    return this.http.post<Post>(`${this.baseUrl}/post/`, {
      title: post.title,
      body: post.body,
    }, { headers: this.headers() });
  }

  public updatePost(post: Post): Observable<Object> {
    if (!post.id || post.id <= 0) {
      return new Observable(o => {
        o.error("Post ID must be a positive integer");
      });
    }

    return this.http.post(`${this.baseUrl}/post/${post.id}`, {
      title: post.title,
      body: post.body,
    }, { headers: this.headers() });
  }

  public deletePost(id: number): Observable<Object> {
    return this.http.delete(`${this.baseUrl}/post/${id}`, { headers: this.headers() });
  }

  public restorePost(id: number): Observable<Object> {
    return this.http.post(`${this.baseUrl}/post/${id}`, {
      deleted: false,
    }, { headers: this.headers() });
  }

  private headers(headers = {}) {
    const token = window.sessionStorage.getItem("token");
    headers["Content-Type"] = "application/json"
    if (token) {
      headers["Authorization"] = token
    }
    return headers
  }
}
