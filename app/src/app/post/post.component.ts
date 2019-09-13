import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { BlogService } from '../blog.service';
import { Post } from '../post';
import { AuthenticationService } from '../authentication.service';

@Component({
  selector: 'app-post',
  templateUrl: './post.component.html',
  styleUrls: ['./post.component.less']
})
export class PostComponent implements OnInit {

  public new: boolean;
  public editing: boolean;
  public id: number;
  public action: string;
  public post: Post = new Post();

  constructor(private blogService: BlogService, private router: Router, private authenticationService: AuthenticationService,  route: ActivatedRoute) {
    this.post = new Post();
    route.data.forEach(d => {
      if (!this.isLoggedIn()) {
        return;
      }
      this.new = d.new;
      this.editing = d.edit;
    });
    route.params.forEach(p => {
      if (this.id != p.id) {
        this.id = p.id;
        this.loadPost();
      }
    });
    this.loadPost();
  }

  ngOnInit() {
  }

  public isLoggedIn(): boolean {
    return this.authenticationService.isLoggedIn();
  }

  public save() {
    if (this.new) {
      this.blogService.newPost(this.post).subscribe(post => {
        this.post = post;
        this.router.navigate(["/post", post.id]);
      }, error => {
          // @TODO: show friendly error in UI
          console.error(error)
      });
    } else if (this.editing) {
      this.blogService.updatePost(this.post).subscribe(() => {
        this.router.navigate(["/post", this.id]);
      }, error => {
          // @TODO: show friendly error in UI
          console.error(error)
      });
    }
  }

  public edit() {
    this.router.navigate(["/post", this.id, "edit"]);
  }

  public delete() {
    this.blogService.deletePost(this.id).subscribe(() => {
      this.router.navigate(["/posts"]);
    }, error => {
      // @TODO: show friendly error in UI
      console.error(error)
    });
  }

  public restore() {
    this.blogService.restorePost(this.id).subscribe(() => {
      this.router.navigate(["/posts"]);
    }, error => {
      // @TODO: show friendly error in UI
      console.error(error)
    });
  }

  private loadPost() {
    if (!this.new) {
      this.blogService.getPost(this.id).subscribe(post => {
        this.post = post;
      }, error => {
          // @TODO: show friendly error in UI
          console.error(error)
      });
    }
  }

}
