import { Component, OnInit } from '@angular/core';
import { BlogService, BlogPostFilter } from '../blog.service';
import { ActivatedRoute } from '@angular/router';
import { Post } from '../post';

@Component({
  selector: 'app-post-list',
  templateUrl: './post-list.component.html',
  styleUrls: ['./post-list.component.less']
})
export class PostListComponent implements OnInit {

  public posts: Post[] = [];
  public filter: BlogPostFilter;
  public deleted: boolean = false;

  private loadAttempted: boolean = false;

  constructor(private blogService: BlogService, private route: ActivatedRoute) {}

  ngOnInit() {
    this.route.data.subscribe(d => {
      if (this.filter !== d.filter) {
        this.filter = d.filter;
        this.deleted = this.filter === BlogPostFilter.Deleted;
        this.loadPosts();
      }
    });
    if (!this.loadAttempted) {
      this.loadPosts();
    }
  }

  private loadPosts() {
    this.loadAttempted = true;
    this.blogService.getPosts(this.filter).subscribe(posts => {
      this.posts = posts;
    }, error => {
        // @TODO: show friendly error in UI
        console.error(error)
    });
  }

}
