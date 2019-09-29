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
  public errors: string[] = [];

  constructor(private blogService: BlogService, private router: Router, private authenticationService: AuthenticationService,  private route: ActivatedRoute) {}

  ngOnInit() {
    this.post = new Post();
    this.route.data.subscribe(d => {
      if (!this.isLoggedIn()) {
        return;
      }
      this.new = d.new;
      this.editing = d.edit;
    });
    this.route.params.subscribe(p => {
      if (this.id != p.id) {
        this.id = p.id;
        this.loadPost();
      }
    });
  }

  public isLoggedIn(): boolean {
    return this.authenticationService.isLoggedIn();
  }

  public save() {
    if (this.new) {
      this.blogService.newPost(this.post).subscribe(post => {
        this.router.navigate(["/post", post.id]);
      }, error => {
        this.showError("Sorry, there was a problem saving this blog post", error)
      });
    } else if (this.editing) {
      this.blogService.updatePost(this.post).subscribe(() => {
        this.router.navigate(["/post", this.id]);
      }, error => {
        this.showError("Sorry, there was a problem saving this blog post", error)
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
      this.showError("Sorry, there was a problem deleting this blog post", error)
    });
  }

  public restore() {
    this.blogService.restorePost(this.id).subscribe(() => {
      this.router.navigate(["/posts"]);
    }, error => {
      this.showError("Sorry, there was a problem restoring this blog post", error)
    });
  }

  private loadPost() {
    if (this.new) {
      return
    }

    this.blogService.getPost(this.id).subscribe(post => {
      if (!post) {
        this.showError("Sorry, there was a problem retrieving this blog post", post)
        return
      }
      this.post = post;
    }, error => {
        this.showError("Sorry, there was a problem retrieving this blog post", error)
    });
  }

  private showError(msg: string, detail: any) {
    this.errors.push(msg)
    if (detail === undefined) {
      console.error(msg)
    } else {
      console.error(msg, ":" , detail)
    }
  }

}
